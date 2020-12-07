package xls

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx/v3"
)

// NewFromTemplate copies a template file located at templatePath to the filename and returns
// it ready for writing
func NewFromTemplate(filename string, templatePath string) (*xlsx.File, error) {
	sourceStat, err := os.Stat(templatePath)
	if err != nil {
		return nil, errors.Wrap(err, "VAT template file read failed")
	}

	if !sourceStat.Mode().IsRegular() {
		return nil, errors.Wrap(err, "VAT template file mode error")
	}

	source, err := os.Open(templatePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open VAT template file")
	}
	defer source.Close()

	dest, err := os.Create(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s for writing", filename)
	}

	if _, err := io.Copy(dest, source); err != nil {
		return nil, errors.Wrap(err, "Copying from VAT template to destination failed")
	}
	dest.Close()

	return xlsx.OpenFile(filename)

}
