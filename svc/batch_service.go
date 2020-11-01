package svc

import (
	"time"

	"github.com/BTBurke/vatinator/db"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type BatchService interface {
	CreateBatch(acctID string, startID string) (*Batch, error)
	GetBatch(acctID string, id string) (*Batch, error)
	CloseBatch(acctID string, id string) error
}

type b struct {
	db *badger.DB
}

func (b b) CreateBatch(acctID string, startID string) (*Batch, error) {
	batch := &Batch{StartID: startID}
	key := &BatchKey{
		AccountID: acctID,
		BatchID:   xid.New().String(),
	}

	if err := b.db.Update(func(txn *badger.Txn) error {
		return db.Set(txn, key, batch)
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create batch")
	}
	return batch, nil
}

func (b b) GetBatch(acctID string, batchID string) (*Batch, error) {
	key := &BatchKey{acctID, batchID}
	batch := &Batch{}

	if err := b.db.View(func(txn *badger.Txn) error {
		if _, err := getBatch(txn, key, batch, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return batch, nil
}

// getBatch will return a materialized view of this batch.  If it is not closed, it will
// get all receipts associated with the open batch to calculate the current value of the batch
// IMPORTANT: setting persistView will save the materialized view and must be in a mutable badger transaction
func getBatch(txn *badger.Txn, key *BatchKey, batch *Batch, persistView bool) ([]Receipt, error) {
	if err := db.Get(txn, key, batch); err != nil {
		return nil, err
	}

	// materialize the view of receipt totals inside the transaction
	receipts, err := getReceiptsForBatch(txn, key)
	if err != nil {
		return nil, err
	}

	batch.NumReceipts = len(receipts)

	total := 0
	vat := 0
	for _, r := range receipts {
		total += r.Total
		vat += r.VAT
	}
	batch.VAT = vat
	batch.Total = total

	if persistView {
		if err := db.Set(txn, key, batch); err != nil {
			return nil, err
		}
	}
	return receipts, nil
}

func (b b) CloseBatch(accountID string, batchID string) error {
	key := &BatchKey{accountID, batchID}
	batch := &Batch{}

	if err := b.db.Update(func(txn *badger.Txn) error {
		if _, err := getBatch(txn, key, batch, false); err != nil {
			return err
		}
		batch.Closed = time.Now().Unix()
		return db.Set(txn, key, batch)
	}); err != nil {
		return err
	}
	return nil
}

var _ BatchService = b{}
