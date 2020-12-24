package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"gopkg.in/yaml.v2"
)

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

func main() {
	for _, r := range []rune("âãÏÓ") {
		fmt.Println(r)
	}

	fields, err := loadFields()
	if err != nil {
		log.Fatal(err)
	}

	data := map[fieldKey]string{
		type1:    "Gas 95",
		content1: "shart",
	}
	fdf := createFDF(data, fields)
	if err := ioutil.WriteFile("data.fdf", fdf, 0644); err != nil {
		log.Fatal(err)
	}

}

func loadFields() (map[fieldKey]string, error) {
	replacements := map[string]string{
		"&#228;": "ä",
		"&#245;": "õ",
	}

	data, err := ioutil.ReadFile("fields.yaml")
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

func createFDF(data map[fieldKey]string, names map[fieldKey]string) []byte {

	e := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	enc := e.NewEncoder()

	// data := map[string]string{
	// 	"Välisesinduse nimi The name of the foreign representation": "Test",
	// }

	b := new(bytes.Buffer)

	b.Write([]byte(fdfHeader))
	for key, value := range data {
		b.Write([]byte("<<\n/T ("))
		k, err := enc.Bytes([]byte(fmt.Sprintf("%s", names[key])))
		if err != nil {
			log.Fatal(err)
		}
		b.Write(k)
		b.Write([]byte(")\n"))
		b.Write([]byte("/V ("))
		v, err := enc.Bytes([]byte(fmt.Sprintf("%s", value)))
		if err != nil {
			log.Fatal(err)
		}
		b.Write(v)
		b.Write([]byte(")\n>>\n"))
	}
	b.Write([]byte(fdfFooter))
	return b.Bytes()
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
