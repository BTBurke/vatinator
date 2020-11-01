package svc

import (
	"fmt"

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
	fmt.Println(receipts)

	return nil
}
