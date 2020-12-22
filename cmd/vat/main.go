package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BTBurke/clt"
	"github.com/BTBurke/vatinator/bundled"
	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/svc"
	"github.com/dgraph-io/badger/v2"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type config struct {
	FirstName    string
	LastName     string
	FullName     string
	DiplomaticID string
	Embassy      string
	Address      string
	Bank         string
}

type task struct {
	name  string
	image img.Image
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)
		os.Exit(0)
	}

	if _, err := os.Stat(".cfg/key.json"); os.IsNotExist(err) {
		if err := decryptKeyFile(); err != nil {
			log.Fatalf("Failed to decrypt key file: %s", err)
		}
	}

	var cfg config
	if _, err := os.Stat(".cfg/config.json"); os.IsNotExist(err) {
		if err := setup(&cfg); err != nil {
			log.Fatalf("Setup failed: %s", err)
		}
	} else {
		data, err := ioutil.ReadFile(".cfg/config.json")
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			log.Fatal(err)
		}
	}

	template, err := bundled.Asset("assets/vat-template.xlsx")
	if err != nil {
		log.Fatal("Failed to find VAT form template")
	}

	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: needs to keep track of batch IDs based on the directory, not reprocess every
	// receipt every time
	accountID := "1"
	batchID := "1"

	// Walk top level directory to find all potential directories with receipts
	var dirs []string
	if err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		switch {
		case info.IsDir() && path != wd && !strings.HasSuffix(path, ".cfg"):
			dirs = append(dirs, path)
			return filepath.SkipDir
		default:
			return nil
		}
	}); err != nil {
		log.Fatalf("Failed to scan directories for receipt images: %s", err)
	}

	paths := make(map[string]string)
	for i, path := range dirs {
		key := strconv.Itoa(i)
		paths[key] = path
	}
	i := clt.NewInteractiveSession()
	dir := i.AskFromTable("Choose directory with this submission's receipts", paths, "")
	i.Reset()
	dirIndex, err := strconv.Atoi(dir)
	if err != nil {
		log.Fatalf("Unknown directory: %s", dir)
	}

	// Walk directory and find all images to process
	rcptFinder := clt.NewProgressSpinner("Finding receipts")
	rcptFinder.Start()

	var tasks []task
	if err := filepath.Walk(dirs[dirIndex], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		image, err := img.NewImageFromReader(f)
		if err != nil {
			return nil
		}
		tasks = append(tasks, task{path, image})

		return nil
	}); err != nil {
		rcptFinder.Fail()
		log.Fatal(err)
	}
	rcptFinder.Success()

	// Figure out what month and year this batch is for
	y, m, d := time.Now().Date()
	switch {
	case m == time.January && d <= 8:
		y -= 1
		m = time.December
	case d <= 8:
		m = time.Month(int(m) - 1)
	default:
	}
	i.Say("Found %d receipts", len(tasks))
	month := i.AskWithDefault("Enter month number for this submission (e.g. 11 for Nov)", strconv.Itoa(int(m)), clt.AllowedOptions([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}))
	monthIndex, err := strconv.Atoi(month)
	if err != nil {
		log.Fatalf("Invalid month: %s", month)
	}
	// if didn't pick the default, make sure year is correct
	if monthIndex != int(m) {
		i.Reset()
		yearS := i.AskWithDefault("Enter year for this submission", strconv.Itoa(y))
		var err error
		y, err = strconv.Atoi(yearS)
		if err != nil {
			log.Fatalf("Invalid year: %s", yearS)
		}
	}

	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	monthString := months[monthIndex-1]

	// Start parallel processor and wait until finished
	outdir := filepath.Join(dirs[dirIndex], "out")
	if err := os.RemoveAll(outdir); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		if err := os.Mkdir(outdir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %s", err)
		}
	}

	it := clt.NewIncrementalProgressBar(len(tasks), "Doing magic to extract data from receipts")
	it.Start()

	proc := svc.NewParallelProcessor(db, accountID, batchID, &svc.ParallelOptions{
		ReprocessOnRulesChange: true,
		NumProcs:               20,
		Hooks: &svc.Hooks{
			AfterEach: func(r *svc.Receipt) error {
				it.Increment()
				if err := svc.WriteErrors(filepath.Join(outdir, "errors.txt"))(r); err != nil {
					return err
				}
				return nil
			},
		},
		KeyPath: ".cfg/key.json",
	})

	for _, task := range tasks {
		if err := proc.Add(task.name, task.image); err != nil {
			log.Fatalf("Failed when processing %s: %s", task.name, err)
		}
	}
	if err := proc.Wait(); err != nil {
		it.Fail()
		log.Fatalf("Processing images failed: %s", err)
	}
	it.Success()

	// export images to PDFs and fill forms
	exp := clt.NewProgressSpinner("Filling out your forms")
	exp.Start()

	export := svc.NewExportService(db)
	if err := export.Create(accountID, batchID, &svc.ExportOptions{
		FirstName:    cfg.FirstName,
		LastName:     cfg.LastName,
		FullName:     cfg.FullName,
		DiplomaticID: cfg.DiplomaticID,
		Month:        monthString,
		MonthInt:     monthIndex,
		Year:         y,
		Stamp:        []string{cfg.FullName, cfg.Embassy, cfg.Address},
		Bank:         cfg.Bank,
		Template:     template,
		OutputDir:    outdir,
	}); err != nil {
		exp.Fail()
		log.Fatal(err)
	}
	exp.Success()

	i.Reset()
	i.Say("Partial success! See output in %s and review the forms and errors.txt to fix my failings.", filepath.Join(dirs[dirIndex], "out"))
	i.Pause()
}

func openDB() (*badger.DB, error) {
	tmpdir, err := ioutil.TempDir("", "vat")
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(tmpdir)
	opts.Logger = nilLogger{}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setup(cfg *config) error {

	i := clt.NewInteractiveSession()

	i.Say("No configuration exists. Let's set this up.\n")

	cfg.FirstName = i.Ask("Enter your first name", clt.Required())
	i.Reset()

	cfg.LastName = i.Ask("\nEnter your last name", clt.Required())
	i.Reset()

	defName := fmt.Sprintf("%s %s", cfg.FirstName, cfg.LastName)
	i.Say("Your full name will appear as " + defName + ". If this is not ok, enter how you want it to appear.  To accept this, just press enter.")
	cfg.FullName = i.AskWithDefault("Enter full name", defName)
	i.Reset()

	cfg.DiplomaticID = i.AskWithHint("\nEnter your diplomatic ID", "Starts with B in upper right of dip ID", clt.Required())
	i.Reset()

	cfg.Embassy = i.AskWithDefault("\nEnter embassy", "US Embassy")
	i.Reset()

	cfg.Address = i.AskWithDefault("\nEnter embassy address", "Kentmanni 20")
	i.Reset()

	i.Say("Now let's set up your banking details.  First choose your bank then input your bank account number.  These numbers never leave your computer.")

	banks := map[string]string{
		"1": "SEB Bank",
		"2": "Swedbank",
		"3": "LHV Bank",
		"4": "Luminor Bank",
		"5": "Some other bank",
	}
	bankAddresses := map[string]string{
		"1": "AS SEB Bank, EEUHEE2X, Tornimäe 2, 15010 Tallinn, Estonia,",
		"2": "Swedbank AS, HABAEE2X, Liivalaia 8, 15040 Tallinn, Estonia,",
		"3": "AS LHV Bank, LHVBEE22, Tartu mnt 2, 10145 Tallinn, Estonia,",
		"4": "Luminor Bank AS, NDEAEE2X, Liivalaia 45, 10145 Tallinn, Estonia,",
	}
	bank := i.AskFromTable("Choose your bank", banks, "")
	switch bank {
	case "oth":
		i.Reset()
		i.Say("You need to enter the full banking details line.  It should look something like:\nSwedbank AS, HABAEE2X, Liivalaia 8, 15040 Tallinn, Estonia, <your account number>")
		cfg.Bank = i.Ask("Enter bank details")
	default:
		i.Reset()
		acct := i.Ask("Enter your account number", clt.Required())
		cfg.Bank = fmt.Sprintf("%s %s", bankAddresses[bank], acct)
	}
	i.Reset()

	i.Say("Your configuration:\n%s (%s), %s, %s\n%s", cfg.FullName, cfg.DiplomaticID, cfg.Embassy, cfg.Address, cfg.Bank)
	ans := i.AskYesNo("Is this ok", "yes")
	if clt.IsNo(ans) {
		return setup(cfg)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(".cfg/config.json", data, 0644); err != nil {
		return err
	}

	return nil
}

func decryptKeyFile() error {
	dataB64, err := bundled.Asset("assets/api.bin")
	if err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(string(dataB64))
	if err != nil {
		return err
	}

	saltB64, err := bundled.Asset("assets/salt.bin")
	if err != nil {
		return err
	}
	salt, err := base64.StdEncoding.DecodeString(string(saltB64))
	if err != nil {
		return err
	}

	if _, err := os.Stat(".cfg"); os.IsNotExist(err) {
		if err := os.Mkdir(".cfg", 0755); err != nil {
			return err
		}
	}

	i := clt.NewInteractiveSession()

	passphrase := i.Say("The API key file has not been decrypted.").Ask("Enter passphrase")
	passphrase = strings.Trim(passphrase, "\n\r ")

	p := clt.NewProgressSpinner("Decrypting key")
	p.Start()

	dk, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return err
	}
	var secretKey [32]byte
	copy(secretKey[:], dk)

	var decryptNonce [24]byte
	copy(decryptNonce[:], data[:24])
	decrypted, ok := secretbox.Open(nil, data[24:], &decryptNonce, &secretKey)
	if !ok {
		p.Fail()
		return fmt.Errorf("passphrase might be wrong")
	}

	if err := ioutil.WriteFile(".cfg/key.json", decrypted, 0644); err != nil {
		p.Fail()
		return err
	}
	p.Success()

	return nil
}

// gets rid of badger logging messages in terminal
type nilLogger struct{}

func (nilLogger) Errorf(s string, v ...interface{})   {}
func (nilLogger) Warningf(s string, v ...interface{}) {}
func (nilLogger) Infof(s string, v ...interface{})    {}
func (nilLogger) Debugf(s string, v ...interface{})   {}

var _ badger.Logger = nilLogger{}
