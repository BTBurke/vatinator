package main

import (
	"encoding/json"
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
	it := clt.NewIncrementalProgressBar(len(tasks), "Doing magic to extract data from receipts")
	it.Start()

	proc := svc.NewParallelProcessor(db, accountID, batchID, &svc.ParallelOptions{
		ReprocessOnRulesChange: true,
		NumProcs:               20,
		Hooks: &svc.Hooks{
			AfterEach: func(r *svc.Receipt) error { it.Increment(); return nil },
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
		OutputDir:    filepath.Join(dirs[dirIndex], "out"),
	}); err != nil {
		exp.Fail()
		log.Fatal(err)
	}
	exp.Success()

	i.Reset()
	i.Say("Partial success! See output in %s and review the forms and errors.txt to fix my failings.", filepath.Join(dirs[dirIndex], "out"))
}

func openDB() (*badger.DB, error) {
	tmpdir, err := ioutil.TempDir("", "vat")
	if err != nil {
		return nil, err
	}

	db, err := badger.Open(badger.DefaultOptions(tmpdir))
	if err != nil {
		return nil, err
	}
	return db, nil
}
