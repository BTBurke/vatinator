package svc

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/pdf"
	"github.com/BTBurke/vatinator/types"
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
	FirstName         string
	LastName          string
	FullName          string
	Bank              string
	DiplomaticID      string
	Month             string
	MonthInt          int
	Year              int
	Embassy           string
	Stamp             []string
	OutputDir         string
	ConvertXLS2PDF    bool
	FillExciseOptions *pdf.FillExciseOptions
	// template for the VAT form in XLS
	Template []byte
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

	if _, err := os.Stat(opts.OutputDir); os.IsNotExist(err) {
		if err := os.Mkdir(opts.OutputDir, 0755); err != nil {
			return err
		}
	}

	batchKey := &BatchKey{
		AccountID: accountID,
		BatchID:   batchID,
	}

	receipts, err := getReceiptsForBatch(txn, batchKey)
	if err != nil {
		return err
	}

	sort.Slice(receipts, func(i, j int) bool {
		return stringToDate(receipts[i].Date).UTC().Before(stringToDate(receipts[j].Date))
	})

	if err := writeInvoices(txn, accountID, receipts, vat, opts); err != nil {
		return err
	}
	if err := writeVATForm(receipts, opts); err != nil {
		return err
	}

	// find excise receipts and fill excise form
	var excises []types.Excise
	var exciseReceipts []Receipt
	for _, r := range receipts {
		if r.IsExcise {
			exciseReceipts = append(exciseReceipts, r)
			excises = append(excises, types.Excise{
				Type:    r.ExciseType,
				Amount:  r.ExciseAmount,
				Arve:    r.ReceiptNumber,
				Content: "", // empty string for gas receipts
				Date:    r.Date,
			})
		}
	}
	if len(exciseReceipts) > 0 {
		if err := writeInvoices(txn, accountID, exciseReceipts, excise, opts); err != nil {
			return err
		}
		if err := writeExciseForm(excises, opts); err != nil {
			return err
		}
	}

	return nil
}

type invoiceType string

const (
	vat    invoiceType = "Invoices"
	excise invoiceType = "Excise"
)

func writeInvoices(txn *badger.Txn, accountID string, receipts []Receipt, t invoiceType, opts *ExportOptions) error {
	var perPacket int
	switch t {
	case vat:
		perPacket = 17
	case excise:
		perPacket = 6
	default:
		return fmt.Errorf("unknown invoice type %s", t)
	}

	packets := len(receipts)/perPacket + 1
	for packet := 0; packet < packets; packet++ {

		var fpath string
		switch t {
		case vat:
			fpath = filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-VAT-%s%d-Invoices%d.pdf", opts.LastName, opts.Month, opts.Year, packet+1))
		case excise:
			fpath = filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-Excise-%s%d-Fuel_Invoices%d.pdf", opts.LastName, opts.Month, opts.Year, packet+1))
		}
		p := pdf.NewPDF(fpath)

		// Write each receipt to as page in PDF
		for i := 0; i < perPacket; i++ {
			current := packet*perPacket + i
			if current >= len(receipts) {
				continue
			}
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
		}
		if err := p.Save(); err != nil {
			return errors.Wrap(err, "failed to save pdf")
		}
	}

	return nil
}

func writeExciseForm(receipts []types.Excise, opts *ExportOptions) error {
	packets := len(receipts)/6 + 1
	for packet := 0; packet < packets; packet++ {
		excisePath := filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-Excise-%s%d-Fuel_Form%d.pdf", opts.LastName, opts.Month, opts.Year, packet+1))
		if err := pdf.FillExcise(excisePath, receipts, types.ExciseMetadata{
			Bank:    opts.Bank,
			Name:    opts.FullName,
			Embassy: opts.Embassy,
			Date:    fmt.Sprintf("%s %d", opts.Month, opts.Year),
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

func writeVATForm(receipts []Receipt, opts *ExportOptions) error {
	packets := len(receipts)/17 + 1
	for packet := 0; packet < packets; packet++ {

		xpath := filepath.Join(opts.OutputDir, fmt.Sprintf("USA-%s-VAT-%s%d-VAT%d.xlsx", opts.LastName, opts.Month, opts.Year, packet+1))
		xlsfile, err := xls.NewFromTemplate(xpath, opts.Template)
		if err != nil {
			return errors.Wrap(err, "failed to create new VAT file")
		}

		// write form header information
		if err := xls.WriteName(opts.FullName, xlsfile); err != nil {
			return err
		}
		if err := xls.WriteDipNumber(opts.DiplomaticID, xlsfile); err != nil {
			return err
		}
		if err := xls.WriteSubmissionMonth(opts.MonthInt, opts.Year, xlsfile); err != nil {
			return err
		}
		if err := xls.WriteBankInfo(opts.Bank, xlsfile); err != nil {
			return err
		}

		// Write each VAT line
		for i := 0; i < 17; i++ {
			current := packet*17 + i
			if current >= len(receipts) {
				continue
			}
			receipt := &receipts[current]

			// write excel line
			if err := xls.WriteVATLine(xlsfile, receipt, i); err != nil {
				return errors.Wrap(err, "failed to write line to VAT file")
			}
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
