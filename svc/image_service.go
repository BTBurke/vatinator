package svc

import (
	"github.com/BTBurke/vatinator/db"
	"github.com/BTBurke/vatinator/img"
	"github.com/dgraph-io/badger/v2"
)

type ImageService interface {
	Upsert(accountID string, receiptID string, i img.Image) error
	Get(accountID string, receiptID string) (img.Image, error)
}

type i struct {
	db *badger.DB
}

func (i i) Upsert(accountID string, receiptID string, image img.Image) error {
	return i.db.Update(func(txn *badger.Txn) error {
		return upsertImage(txn, accountID, receiptID, image)
	})
}

func upsertImage(txn *badger.Txn, accountID string, receiptID string, i img.Image) error {
	key := &ImageKey{
		AccountID: accountID,
		ReceiptID: receiptID,
	}
	return db.Set(txn, key, &i)
}

func (i i) Get(accountID string, receiptID string) (img.Image, error) {

	var image img.Image
	if err := i.db.View(func(txn *badger.Txn) error {
		var err error
		image, err = getImage(txn, accountID, receiptID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return image, err
	}
	return image, nil
}

func getImage(txn *badger.Txn, accountID string, receiptID string) (img.Image, error) {
	key := &ImageKey{
		AccountID: accountID,
		ReceiptID: receiptID,
	}
	image := img.Image{}
	if err := db.Get(txn, key, &image); err != nil {
		return image, err
	}
	return image, nil
}

func NewImageService(db *badger.DB) ImageService {
	return i{db}
}

var _ ImageService = i{}
