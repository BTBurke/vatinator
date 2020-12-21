package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BTBurke/clt"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

// Encrypts API keys with secretbox using a passphrase that is turned into a key using scrypt and a random salt.
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: encrypt <filename>")
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	passphrase := clt.NewInteractiveSession().Ask("Passphrase")
	passphrase = strings.Trim(passphrase, "\n")

	var salt [8]byte
	if _, err := io.ReadFull(rand.Reader, salt[:]); err != nil {
		log.Fatal(err)
	}

	sOut := base64.StdEncoding.EncodeToString(salt[:])

	if err := ioutil.WriteFile("./assets/salt.bin", []byte(sOut), 0644); err != nil {
		log.Fatal(err)
	}

	dk, err := scrypt.Key([]byte(passphrase), salt[:], 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	encrypted, err := encrypt(dk, data)
	if err != nil {
		log.Fatal(err)
	}

	eOut := base64.StdEncoding.EncodeToString(encrypted)

	if err := ioutil.WriteFile("./assets/api.bin", []byte(eOut), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encrypted files done.  Don't forget to run go-bindata to embed them in binary.")
	os.Exit(0)

}

// uses NaCl secretbox to encrypt data
func encrypt(key []byte, data []byte) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], key)

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	encrypted := secretbox.Seal(nonce[:], data, &nonce, &secretKey)

	return encrypted, nil
}
