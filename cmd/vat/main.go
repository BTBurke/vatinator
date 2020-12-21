package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BTBurke/clt"
	"github.com/BTBurke/vatinator/bundled"
	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/svc"
	"github.com/dgraph-io/badger/v2"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const help = `
Usage:
	vat <directory of photos>
`

func main() {

	if _, err := os.Stat(".cfg/key.json"); os.IsNotExist(err) {
		if err := decryptKeyFile(); err != nil {
			log.Fatalf("Failed to decrypt key file: %s", err)
		}
	}
	log.Fatal("stopping")

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

	log.Printf("salt length:  %d", len(salt))

	if _, err := os.Stat(".cfg"); os.IsNotExist(err) {
		if err := os.Mkdir(".cfg", 0755); err != nil {
			return err
		}
	}

	i := clt.NewInteractiveSession()

	passphrase := i.Say("The API key file has not been decrypted.").Ask("Enter passphrase")
	passphrase = strings.Trim(passphrase, "\n")

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
		return fmt.Errorf("passphrase might be wrong")
	}

	if err := ioutil.WriteFile(".cfg/key.json", decrypted, 0644); err != nil {
		return err
	}

	return nil
}
