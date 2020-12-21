package svc

import (
	"fmt"
	"os"
)

// Hooks allow arbitrary hooks to run before/after processing a batch or before/after each receipt
type Hooks struct {
	// Runs before the batch is processed
	BeforeStart func()
	// Runs after all receipts are processed
	AfterEnd func()
	// Runs before each receipt is processed
	BeforeEach func(r *Receipt) error
	// Runs after a receipt is processed
	AfterEach func(r *Receipt) error
}

type ReceiptHook func(r *Receipt) error

// WriteErrors writes errors to a file after each receipt is processed
func WriteErrors(file string) ReceiptHook {
	return func(r *Receipt) error {
		var f *os.File
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err = os.Create(file)
			if err != nil {
				return err
			}
		} else {
			f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				return err
			}
		}
		defer f.Close()

		for _, e := range r.Errors {
			if _, err := f.Write([]byte(fmt.Sprintf("%s: %s\n", r.Filename, e))); err != nil {
				return err
			}
		}
		return nil
	}
}
