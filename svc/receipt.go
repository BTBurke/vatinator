package svc

import (
	"fmt"
	"time"

	vat "github.com/BTBurke/vatinator"
	"github.com/BTBurke/vatinator/db"
	"github.com/vmihailenco/msgpack/v5"
)

// keep receipts for 95 days
var ReceiptDuration time.Duration = 24 * time.Hour * 95

type Receipt struct {
	ID            string
	Vendor        string
	TaxID         string
	Total         int
	VAT           int
	Date          string
	ReceiptNumber string
	BatchID       string
}

func (r *Receipt) Type() byte {
	return db.Receipt
}

func (r *Receipt) TTL() time.Duration {
	return ReceiptDuration
}

func (r *Receipt) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(r)
}

func (r *Receipt) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, r)
}

func (r *Receipt) GetVendor() string {
	return r.Vendor
}

func (r *Receipt) GetReceiptNumber() string {
	return r.ReceiptNumber
}

func (r *Receipt) GetTaxID() string {
	return r.TaxID
}

func (r *Receipt) GetDate() string {
	return r.Date
}

func (r *Receipt) GetTotal() int {
	return r.Total
}

func (r *Receipt) GetVAT() int {
	return r.VAT
}

type ReceiptKey struct {
	AccountID string
	ReceiptID string
}

func (rk *ReceiptKey) MarshalBinary() ([]byte, error) {
	if len(rk.AccountID) == 0 || len(rk.ReceiptID) == 0 {
		return nil, fmt.Errorf("receipt key error: acct: %s receipt: %s", rk.AccountID, rk.ReceiptID)
	}
	return []byte(fmt.Sprintf("a/%s/r/%s", rk.AccountID, rk.ReceiptID)), nil
}

func (rk *ReceiptKey) UnmarshalBinary(data []byte) error {
	key := splitKey(data)
	acctID, ok := key["a"]
	if !ok {
		return fmt.Errorf("receipt missing account ID: %s", string(data))
	}
	receiptID, ok := key["r"]
	if !ok {
		return fmt.Errorf("receipt missing receipt ID: %s", string(data))
	}

	rk.AccountID = acctID
	rk.ReceiptID = receiptID

	return nil
}

var _ db.Entity = &Receipt{}
var _ db.Key = &ReceiptKey{}
var _ vat.VATLine = &Receipt{}
