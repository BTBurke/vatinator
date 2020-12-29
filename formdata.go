package vatinator

import (
	"encoding/json"
)

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

func UnmarshalFormData(b []byte) (FormData, error) {
	var fd FormData
	if err := json.Unmarshal(b, &fd); err != nil {
		return FormData{}, err
	}
	return fd, nil
}

func required(fields ...string) bool {
	for _, field := range fields {
		if len(field) == 0 {
			return false
		}
	}
	return true
}
