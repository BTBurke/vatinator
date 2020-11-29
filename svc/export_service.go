package svc

import (
	"fmt"
	"log"

	"github.com/BTBurke/vatinator/pdf"
	"github.com/dgraph-io/badger/v2"
)

type ExportService interface {
	Create(accountID string, batchID string, options *ExportOptions) error
	Get(accountID string, batchID string) (*Export, error)
	Del(accountID string, batchID string) error
}

type e struct {
	db *badger.DB
}

func (e e) Create(accountID string, batchID string, options *ExportOptions) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return create(txn, accountID, batchID, options)
	})
}

func create(txn *badger.Txn, accountID string, batchID string, options *ExportOptions) error {
	batchKey := &BatchKey{
		AccountID: accountID,
		BatchID:   batchID,
	}
	batch := &Batch{}

	receipts, err := getBatch(txn, batchKey, batch, false)
	if err != nil {
		return err
	}

	// TODO: sort receipts, create temp dir, populate with export, zip, store
	packets := len(receipts)/17 + 1
	for packet := 0; packet < packets; packet++ {
		p := pdf.NewPDF(fmt.Sprintf("test_%d.pdf", packet))
		for i := 0; i < 17; i++ {
			current := packet*17 + i
			if current >= len(receipts) {
				continue
			}
			id := receipts[current].ID
			img, err := getImage(txn, accountID, id)

			if err != nil {
				log.Fatalf("failed to get image: %v", err)
			}
			if err := p.WriteReceipt(img); err != nil {
				log.Fatalf("failed to write receipt to pdf: %v", err)
			}
		}
		if err := p.Save(); err != nil {
			log.Fatalf("failed to save pdf: %v", err)
		}
	}

	return nil
}
