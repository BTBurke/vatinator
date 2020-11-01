package svc

import (
	"fmt"

	"github.com/BTBurke/vatinator/db"
	"github.com/dgraph-io/badger/v2"
)

type ReceiptService interface {
	Upsert(accountID string, r *Receipt) error
	Get(accountID string, receiptID string) (*Receipt, error)
	Del(accountID string, receiptID string) error
	GetBatch(accountID string, batchID string) ([]Receipt, error)
}

type r struct {
	db *badger.DB
}

func (r r) Upsert(accountID string, receipt *Receipt) error {
	rk := &ReceiptKey{accountID, receipt.ID}
	return r.db.Update(func(txn *badger.Txn) error {
		return db.Set(txn, rk, receipt)
	})
}

func (r r) Get(accountID string, receiptID string) (*Receipt, error) {
	receipt := &Receipt{}
	key := &ReceiptKey{accountID, receiptID}

	if err := r.db.View(func(txn *badger.Txn) error {
		return getReceipt(txn, key, receipt)
	}); err != nil {
		return nil, err
	}
	return receipt, nil
}

func getReceipt(txn *badger.Txn, key *ReceiptKey, receipt *Receipt) error {
	return db.Get(txn, key, receipt)
}

func (r r) Del(accountID string, receiptID string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return db.Del(txn, &ReceiptKey{accountID, receiptID})
	})
}

func (r r) GetBatch(accountID string, batchID string) ([]Receipt, error) {
	var receipts []Receipt
	if err := r.db.View(func(txn *badger.Txn) error {
		rcpts, err := getReceiptsForBatch(txn, &BatchKey{accountID, batchID})
		if err != nil {
			return err
		}
		receipts = append(receipts, rcpts...)
		return nil
	}); err != nil {
		return nil, err
	}
	return receipts, nil
}

func NewReceiptService(db *badger.DB) ReceiptService {
	return r{db}
}

// returns all receipts for a batch
// TODO: currently reverse iterates through all receipts associated with an account, might be a better
// way to do this by recording start and ending receipt ID per batch or something
func getReceiptsForBatch(txn *badger.Txn, key *BatchKey) ([]Receipt, error) {
	account := key.AccountID
	batch := key.BatchID

	if len(batch) == 0 {
		return nil, fmt.Errorf("failed to get receipts, nil batch ID")
	}
	opt := badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   25,
		Reverse:        true,
		AllVersions:    false,
	}
	it := txn.NewIterator(opt)
	defer it.Close()

	prefix := iterateReceipt(account)
	start := iterateReceiptEnd(account)
	var errs []error

	var receipts []Receipt
	for it.Seek(start); it.ValidForPrefix(prefix); it.Next() {
		var r *Receipt
		if err := db.FromItem(it.Item(), r); err != nil {
			errs = append(errs, err)
			continue
		}
		if r.BatchID == batch {
			receipts = append(receipts, *r)
		}
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("errors getting receipts for batch %s: %v", batch, errs)
	}
	return receipts, nil
}

// returns a starting point in the tree of the form a/[account id/r/ which should start iteration on
// on all receipts for this account
func iterateReceipt(accountID string) []byte {
	return []byte(fmt.Sprintf("a/%s/r/", accountID))
}

// returns an ending point at the end of all receipts for this account for reverse iteration through
// the receipts
func iterateReceiptEnd(accountID string) []byte {
	return append(iterateReceipt(accountID), 0xFF)
}

var _ ReceiptService = r{}
