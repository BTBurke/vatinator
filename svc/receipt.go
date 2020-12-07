package svc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/BTBurke/vatinator/db"
	"github.com/BTBurke/vatinator/xls"
)

// keep receipts for 95 days
var ReceiptDuration time.Duration = 24 * time.Hour * 95

type Precision int

const (
	// Currecy of the form 23,33
	Digit2 Precision = iota + 2
	// Currency of the form 23,333 which is allowed by VAT regulations
	Digit3
)

type Receipt struct {
	ID     string
	Vendor string
	TaxID  string
	Total  int
	VAT    int
	// TODO: Switch to unix time at midnight UTC on day receipt was issued
	Date          string
	ReceiptNumber string

	BatchID string
	// Unix time that the receipt was verified
	Reviewed int64
	// Precision of the currency, 2 or 3 digits
	CurrencyPrecision Precision
	Errors            []string
	// RulesVersion indicates which version of the rules engine was used to process the receipt. It
	// can be used to reprocess the receipt when upgrades to the rules engine are made.  See
	// Processor for options to force recomputation.  Default is reprocessing when rules change.
	RulesVersion string
}

func (r *Receipt) Type() byte {
	return db.Receipt
}

func (r *Receipt) TTL() time.Duration {
	return ReceiptDuration
}

func (r *Receipt) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Receipt) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
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

func (r *Receipt) GetTotal() string {
	switch r.CurrencyPrecision {
	case Digit3:
		return currency3ToString(r.Total)
	default:
		return currency2ToString(r.Total)
	}
}

func (r *Receipt) GetVAT() string {
	switch r.CurrencyPrecision {
	case Digit3:
		return currency3ToString(r.VAT)
	default:
		return currency2ToString(r.VAT)
	}
}

func currency2ToString(d int) string {
	switch {
	case d < 10:
		return fmt.Sprintf("0.0%d", d)
	case d >= 10 && d < 100:
		return fmt.Sprintf("0.%d", d)
	default:
		ds := strconv.Itoa(d)
		return fmt.Sprintf("%s.%s", ds[0:len(ds)-2], ds[len(ds)-2:])
	}
}

func currency3ToString(d int) string {
	switch {
	case d < 10:
		return fmt.Sprintf("0.00%d", d)
	case d >= 10 && d < 100:
		return fmt.Sprintf("0.0%d", d)
	case d >= 100 && d < 1000:
		return fmt.Sprintf("0.%d", d)
	default:
		ds := strconv.Itoa(d)
		return fmt.Sprintf("%s.%s", ds[0:len(ds)-3], ds[len(ds)-3:])
	}
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
var _ xls.VATLine = &Receipt{}
