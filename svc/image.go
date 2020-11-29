package svc

import (
	"fmt"

	"github.com/BTBurke/vatinator/db"
)

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
