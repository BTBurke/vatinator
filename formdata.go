package vat

import (
	"errors"
	"fmt"
	"strconv"
)

// AccountID is the ID in the database
type AccountID int

// Fulfills stringer interface
func (a AccountID) String() string {
	return strconv.Itoa(a)
}

// Valuer interface for SQL
func (a AccountID) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	return int(a), nil
}

// Scanner interface for SQL
func (a *AccountID) Scan(value interface{}) error {
	// if value is nil, error
	if value == nil {
		return fmt.Errorf("accountID is nil")
	}
	bv, err := driver.Int.ConvertValue(value)
	if err == nil {
		// if this is a bool type
		if v, ok := bv.(int); ok {
			// set the value of the pointer yne to YesNoEnum(v)
			*a = int(v)
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan account ID")
}

// FormData is the data needed to fill out the form top sections
type FormData struct {
	// Personal data of the submitter
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	DiplomaticID string `json:"diplomatic_id"`
	// Embassy details
	Embassy string `json:"embassy"`
	Address string `json:"address"`
	// Bank for deposit
	Bank     string `json:"bank"`
	BankName string `json:"bank_name"`
	Account  string `json:"account"`
}

// IsValid checks the form data to ensure that all required fields are set
func (fd FormData) IsValid() bool {
	return required(fd.FirstName, fd.LastName, fd.FullName, fd.DiplomaticID, fd.Embassy, fd.Address, fd.Bank)
}

func required(fields ...string) bool {
	for _, field := range fields {
		if len(field) == 0 {
			return false
		}
	}
	return true
}
