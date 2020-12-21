package xls

import (
	"os"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx/v3"
)

// NewFromTemplate copies a template file located at templatePath to the filename and returns
// it ready for writing
func NewFromTemplate(filename string, template []byte) (*xlsx.File, error) {
	dest, err := os.Create(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s for writing", filename)
	}

	if _, err := dest.Write(template); err != nil {
		return nil, errors.Wrap(err, "Writing VAT template to destination failed")
	}
	dest.Close()

	return xlsx.OpenFile(filename)

}
