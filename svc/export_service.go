package svc

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/pdf"
	"github.com/BTBurke/vatinator/xls"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
)

type ExportService interface {
	Create(accountID string, batchID string, options *ExportOptions) error
	//Get(accountID string, batchID string) (*Export, error)
	//Del(accountID string, batchID string) error
}

type ExportOptions struct {
	FirstName      string
	LastName       string
	Month          string
	Year           int
	Stamp          []string
	OutputDir      string
	ConvertXLS2PDF bool
}

func DefaultExportOptions() *ExportOptions {
	year, month, _ := time.Now().Date()
	return &ExportOptions{
		FirstName: "First",
		LastName:  "Last",
		Month:     strings.ToUpper(month.String()),
		Year:      year,
		OutputDir: "",
	}
}

type e struct {
	db *badger.DB
}

func NewExportService(db *badger.DB) ExportService {
	return e{db}
}

func (e e) Create(accountID string, batchID string, options *ExportOptions) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return create(txn, accountID, batchID, options)
	})
}

func create(txn *badger.Txn, accountID string, batchID string, opts *ExportOptions) error {
	if opts == nil {
		opts = DefaultExportOptions()
	}
	batchKey := &BatchKey{
		AccountID: accountID,
		BatchID:   batchID,
	}

	receipts, err := getReceiptsForBatch(txn, batchKey)
	if err != nil {
		return err
	}
	log.Printf("Found %d receipts...", len(receipts))

	sort.Slice(receipts, func(i, j int) bool {
		return stringToDate(receipts[i].Date).UTC().Before(stringToDate(receipts[j].Date))
	})

	// TODO: sort receipts, create temp dir, populate with export, zip, store
	packets := len(receipts)/17 + 1
	for packet := 0; packet < packets; packet++ {

		p := pdf.NewPDF(filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-VAT-%s%d-Invoices%d.pdf", opts.LastName, opts.Month, opts.Year, packet+1)))

		xpath := filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-VAT-%s%d-VAT%d.xlsx", opts.LastName, opts.Month, opts.Year, packet+1))
		xlsfile, err := xls.NewFromTemplate(xpath, "./vat-template.xlsx")
		if err != nil {
			return errors.Wrap(err, "failed to create new VAT file")
		}

		for i := 0; i < 17; i++ {
			current := packet*17 + i
			if current >= len(receipts) {
				continue
			}
			receipt := &receipts[current]
			id := receipts[current].ID
			image, err := getImage(txn, accountID, id)
			if err != nil {
				return errors.Wrap(err, "failed to get image")
			}

			composited, err := img.CompositeReceipt(i+1, image, opts.Stamp, 0)
			if err != nil {
				return errors.Wrap(err, "failed to composite image")
			}

			if err := p.WriteReceipt(composited); err != nil {
				return errors.Wrap(err, "failed to write receipt to pdf")
			}

			// write excel line
			if err := xls.WriteVATLine(xlsfile, receipt, i); err != nil {
				return errors.Wrap(err, "failed to write line to VAT file")
			}
		}
		if err := p.Save(); err != nil {
			return errors.Wrap(err, "failed to save pdf")
		}
		if err := xlsfile.Save(xpath); err != nil {
			return errors.Wrapf(err, "failed to save VAT file to %s", xpath)
		}
	}

	return nil
}

func stringToDate(d string) time.Time {
	t, err := time.Parse("02/01/2006", d)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}
