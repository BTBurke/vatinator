package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BTBurke/vatinator/svc"
	"github.com/dgraph-io/badger/v2"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmpdir := os.TempDir()
	log.Printf("database directory: %s", tmpdir)
	db, err := badger.Open(badger.DefaultOptions(tmpdir))
	if err != nil {
		log.Fatal(err)
	}

	accountID := "1"
	batchID := "1"

	p := filepath.Join(wd, "nov")
	if err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		log.Printf("Processing %s", path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		proc := svc.NewSingleProcessor(db, accountID, batchID)
		if err := proc.Add(path, f); err != nil {
			log.Fatal(err)
		}
		_ = proc.Wait()

		return nil
	}); err != nil {
		log.Fatal(err)
	}

}
