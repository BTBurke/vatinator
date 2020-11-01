package svc

import (
	"fmt"
	"time"

	"github.com/BTBurke/vatinator/db"
)

var DefaultImageDuration time.Duration = 24 * time.Hour * 95

type Image []byte

func (i *Image) TTL() time.Duration {
	return DefaultImageDuration
}

func (i *Image) Type() byte {
	return db.Image
}

func (i *Image) MarshalBinary() ([]byte, error) {
	return []byte(*i), nil
}

func (i *Image) UnmarshalBinary(data []byte) error {
	copy(*i, data)
	return nil
}

type ImageKey struct {
	AccountID string
	ReceiptID string
}

func (k *ImageKey) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("a/%s/i/%s", k.AccountID, k.ReceiptID)), nil
}

func (k *ImageKey) UnmarshalBinary(data []byte) error {
	key := splitKey(data)
	acctID, ok := key["a"]
	if !ok {
		return fmt.Errorf("image key missing account ID: %s", string(data))
	}
	receiptID, ok := key["i"]
	if !ok {
		return fmt.Errorf("image key missing receipt ID: %s", string(data))
	}

	k.AccountID = acctID
	k.ReceiptID = receiptID

	return nil
}

var _ db.Key = &ImageKey{}
var _ db.Entity = &Image{}
