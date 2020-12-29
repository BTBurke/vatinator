package vatinator

import "github.com/BTBurke/vatinator/svc"

type Options struct {
	CredentialPath string
	ZipResult      bool
	OutputPath     string
	Hooks          svc.Hooks
	FormData       FormData
}

// Process will read receipts located at path and process them into VAT and excise forms.  Returned
// values are full paths to every file generated.
func Process(path string) (string, error) {
	return "", nil
}
