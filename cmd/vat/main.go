package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/svc"
	"github.com/dgraph-io/badger/v2"
)

const help = `
Usage:
	vat <directory of photos>
`

func main() {
	if len(os.Args) != 2 {
		fmt.Print(help)
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmpdir, err := ioutil.TempDir("", "vat")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("database directory: %s", tmpdir)
	db, err := badger.Open(badger.DefaultOptions(tmpdir))
	if err != nil {
		log.Fatal(err)
	}

	export := svc.NewExportService(db)

	accountID := "1"
	batchID := "1"

	p := filepath.Join(wd, os.Args[1])
	proc := svc.NewParallelProcessor(db, accountID, batchID, nil)
	if err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		log.Printf("Processing %s", path)

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		image, err := img.NewImageFromReader(f)
		if err != nil {
			log.Printf("Not an image file, skipping %s", path)
			return nil
		}

		if err := proc.Add(path, image); err != nil {
			log.Fatal(err)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	if err := proc.Wait(); err != nil {
		log.Fatal(err)
	}

	if err := export.Create(accountID, batchID, &svc.ExportOptions{
		FirstName: "Bryan",
		LastName:  "Burke",
		Month:     "November",
		Year:      2020,
		Stamp:     []string{"Bryan Burke", "US Embassy", "Kentmanni 20"},
		OutputDir: p,
	}); err != nil {
		log.Fatal(err)
	}

}
