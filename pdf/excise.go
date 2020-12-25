package pdf

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BTBurke/vatinator/bundled"
	"github.com/BTBurke/vatinator/types"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/yaml.v2"
)

// DefaultURL is the filler service running on gcloud
const DefaultURL string = "https://filler-nd4nplab7a-lz.a.run.app"

type fieldKey string

const (
	embassy  fieldKey = "embassy"
	name              = "name"
	bank              = "bank"
	date              = "date"
	type1             = "type1"
	content1          = "content1"
	amount1           = "amount1"
	excise1           = "excise1"
	arve1             = "arve1"
	type2             = "type2"
	content2          = "content2"
	amount2           = "amount2"
	excise2           = "excise2"
	arve2             = "arve2"
	type3             = "type3"
	content3          = "content3"
	amount3           = "amount3"
	excise3           = "excise3"
	arve3             = "arve3"
	type4             = "type4"
	content4          = "content4"
	amount4           = "amount4"
	excise4           = "excise4"
	arve4             = "arve4"
	type5             = "type5"
	content5          = "content5"
	amount5           = "amount5"
	excise5           = "excise5"
	arve5             = "arve5"
	type6             = "type6"
	content6          = "content6"
	amount6           = "amount6"
	excise6           = "excise6"
	arve6             = "arve6"
	total             = "total"
)

type field struct {
	FieldType          string `yaml:"FieldType"`
	FieldName          string `yaml:"FieldName"`
	FieldNameAlt       string `yaml:"FieldNameAlt"`
	FieldFlags         int    `yaml:"FieldFlags"`
	FieldJustification string `yaml:"FieldJustification"`
	// Short name for looking up the longer field name
	Key string `yaml:"Key"`
}

type FillExciseOptions struct {
	ForceRemote   bool
	DisableRemote bool
	RemoteURL     string
	APIKey        string
}

func FillExcise(path string, rcpts []types.Excise, md types.ExciseMetadata, opts *FillExciseOptions) error {
	if opts == nil {
		opts = &FillExciseOptions{
			DisableRemote: false,
			ForceRemote:   true,
			RemoteURL:     DefaultURL,
			APIKey:        DefaultAPIKey,
		}
	}

	// populate data for form
	data := map[fieldKey]string{
		embassy: md.Embassy,
		name:    md.Name,
		bank:    md.Bank,
		date:    md.Date,
	}

	var tot int
	for i, r := range rcpts {
		m := r.AsMap(i + 1)
		for k, v := range m {
			data[fieldKey(k)] = v
		}
		tot += r.Tax
	}
	data[total] = types.Currency(tot).String()

	// put input files in a temp directory
	tmpdir, err := ioutil.TempDir("", "excise")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	fdfPath := filepath.Join(tmpdir, "data.fdf")
	fdf, err := createFDF(data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fdfPath, fdf, 0644); err != nil {
		log.Fatal(err)
	}

	templatePath := filepath.Join(tmpdir, "template.pdf")
	t, err := bundled.Asset("assets/excise.pdf")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(templatePath, t, 0644); err != nil {
		return err
	}

	// shell out for pdftk to fill form and place at path
	if err := callPdftk(templatePath, fdfPath, path, opts); err != nil {
		return err
	}

	return nil

}

var DefaultAPIKey string = apiKey(".cfg/key.json")

func apiKey(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func loadFields() (map[fieldKey]string, error) {
	replacements := map[string]string{
		"&#228;": "ä",
		"&#245;": "õ",
	}

	data, err := bundled.Asset("assets/fields.yaml")
	if err != nil {
		return nil, err
	}

	out := field{}
	fieldMap := make(map[fieldKey]string)
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		if err := decoder.Decode(&out); err != nil {
			break
		}
		name := out.FieldName
		for old, new := range replacements {
			name = strings.ReplaceAll(name, old, new)
		}
		fieldMap[fieldKey(out.Key)] = name
	}
	return fieldMap, nil
}

func createFDF(data map[fieldKey]string) ([]byte, error) {
	names, err := loadFields()
	if err != nil {
		return nil, err
	}

	e := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	enc := e.NewEncoder()
	b := new(bytes.Buffer)

	b.Write([]byte(fdfHeader))
	for key, value := range data {
		b.Write([]byte("<<\n/T ("))
		k, err := enc.Bytes([]byte(names[key]))
		if err != nil {
			return nil, err
		}
		b.Write(k)
		b.Write([]byte(")\n"))
		b.Write([]byte("/V ("))
		v, err := enc.Bytes([]byte(value))
		if err != nil {
			return nil, err
		}
		b.Write(v)
		b.Write([]byte(")\n>>\n"))
	}
	b.Write([]byte(fdfFooter))
	return b.Bytes(), nil
}

func callPdftk(template string, fdf string, out string, opts *FillExciseOptions) error {
	bin, err := exec.LookPath("pdftk")
	if err != nil || opts.ForceRemote {
		if !opts.DisableRemote && len(opts.APIKey) > 0 && len(opts.RemoteURL) > 0 {
			// no local pdftk so use remote service
			return callPdftkRemote(opts.RemoteURL, fdf, out, opts.APIKey)
		} else {
			return errors.Wrap(err, "no local pdftk and either remote disabled or not API information to call remote service")
		}
	}

	cmd := exec.Command(bin, template, "fill_form", fdf, "output", out)
	stdouterr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(stdouterr))
	}

	return nil
}

func callPdftkRemote(url string, fdf string, out string, apikey string) error {
	if url == "" {
		url = "localhost:8080"
	}
	f, err := os.Open(fdf)
	if err != nil {
		return err
	}
	defer f.Close()

	req, err := http.NewRequest("POST", url, f)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", apikey)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("got status %d: %s", resp.StatusCode, string(data))
	}

	if err := ioutil.WriteFile(out, data, 0644); err != nil {
		return err
	}
	return nil
}

const fdfHeader = `%FDF-1.2
âãÏÓ
1 0 obj 
<<
/FDF 
<<
/Fields [
`

const fdfFooter = `]
>>
>>
endobj 
trailer

<<
/Root 1 0 R
>>
%%EOF
`
