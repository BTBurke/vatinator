package svc

import (
	"github.com/BTBurke/vatinator/db"
	"github.com/dgraph-io/badger/v2"
)

type ImageService interface {
	Upsert(accountID string, receiptID string, image []byte) error
	Get(accountID string, receiptID string) ([]byte, error)
}

type i struct {
	db *badger.DB
}

func (i i) Upsert(accountID string, receiptID string, image []byte) error {
	return i.db.Update(func(txn *badger.Txn) error {
		return upsertImage(txn, accountID, receiptID, image)
	})
}

func upsertImage(txn *badger.Txn, accountID string, receiptID string, image []byte) error {
	key := &ImageKey{
		AccountID: accountID,
		ReceiptID: receiptID,
	}
	img := Image(image)
	return db.Set(txn, key, &img)
}

func (i i) Get(accountID string, receiptID string) ([]byte, error) {

	var out []byte
	if err := i.db.View(func(txn *badger.Txn) error {
		img, err := getImage(txn, accountID, receiptID)
		if err != nil {
			return err
		}
		copy(out, img)
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func getImage(txn *badger.Txn, accountID string, receiptID string) ([]byte, error) {
	key := &ImageKey{
		AccountID: accountID,
		ReceiptID: receiptID,
	}
	img := Image{}
	if err := db.Get(txn, key, &img); err != nil {
		return nil, err
	}
	return []byte(img), nil
}

func NewImageService(db *badger.DB) ImageService {
	return i{db}
}

var _ ImageService = i{}
