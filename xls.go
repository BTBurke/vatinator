package vat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

// number of rows to the first receipt row (0-based)
const rowOffset int = 21

type cellOp func() error

// WriteName writes the applicant's name to both the applicant name field and receiver of refund
func WriteName(name string, f *xlsx.File) error {
	sh := f.Sheets[0]
	if err := setString(8, 0, name, sh); err != nil {
		return fmt.Errorf("error writing name to template: %s", err)
	}
	if err := setString(10, 0, name, sh); err != nil {
		return fmt.Errorf("error writing name to template: %s", err)
	}
	return nil
}

// WriteDipNumber
func WriteDipNumber(num string, f *xlsx.File) error {
	sh := f.Sheets[0]
	if err := setString(8, 3, num, sh); err != nil {
		return fmt.Errorf("error writing dip number to template: %s", err)
	}
	return nil
}

// WriteBankInfo
func WriteBankInfo(info string, f *xlsx.File) error {
	sh := f.Sheets[0]
	if err := setString(12, 0, info, sh); err != nil {
		return fmt.Errorf("error writing bank info to template: %s", err)
	}
	return nil
}

// WriteSubmissionMonth
func WriteSubmissionMonth(month int, year int, f *xlsx.File) error {
	sh := f.Sheets[0]
	if err := setString(16, 0, fmt.Sprintf("%d/01/%d", month, year), sh); err != nil {
		return fmt.Errorf("error writing submission month to spreadsheet: %s", err)
	}
	return nil
}

// WriteVATLine writes a VAT line to the Excel spreadsheet
func WriteVATLine(f *xlsx.File, r *Receipt, num int) error {
	if num > 17 {
		return fmt.Errorf("unallowed row %d: greater than 17", num)
	}

	row := num + rowOffset

	if len(f.Sheets) == 0 {
		return fmt.Errorf("no sheets in file")
	}
	sh := f.Sheets[0]
	ops := []cellOp{
		setNumF(row, 0, num, sh),
		setStringF(row, 1, r.Vendor, sh),
		setStringF(row, 2, r.ID, sh),
		setStringF(row, 3, r.Date, sh),
		setCurrencyF(row, 4, r.Total, sh),
		setCurrencyF(row, 5, r.Tax, sh),
	}

	var errs []string
	for _, op := range ops {
		if err := op(); err != nil {
			errs = append(errs, fmt.Sprintf("%s", err))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func setString(row, col int, s string, sh *xlsx.Sheet) error {
	c, err := sh.Cell(row, col)
	if err != nil {
		return err
	}
	if len(s) == 0 {
		s = "???"
		style := xlsx.NewStyle()
		style.Fill.FgColor = "FF0000FF"
		c.SetStyle(style)
	}
	c.SetString(s)
	return nil
}

// setCurrency takes a single unit currency and converts it to a float without loss of precision
func setCurrency(row, col int, d int, sh *xlsx.Sheet) error {
	c, err := sh.Cell(row, col)
	if err != nil {
		return err
	}
	ds := strconv.Itoa(d)
	switch {
	case d < 10:
		c.SetNumeric(fmt.Sprintf("0.0%d", d))
	case d >= 10 && d < 100:
		c.SetNumeric(fmt.Sprintf("0.%d", d))
	default:
		c.SetNumeric(fmt.Sprintf("%s.%s", ds[0:len(ds)-2], ds[len(ds)-2:]))
	}
	if d == 0 {
		style := xlsx.NewStyle()
		style.Fill.FgColor = "FF0000FF"
		c.SetStyle(style)
	}
	return nil
}

func setNum(row, col int, d int, sh *xlsx.Sheet) error {
	c, err := sh.Cell(row, col)
	if err != nil {
		return err
	}
	c.SetInt(d)
	return nil
}

//
// closures to flatMap over all ops
//
func setNumF(row, col int, d int, sh *xlsx.Sheet) func() error {
	return func() error {
		return setNum(row, col, d, sh)
	}
}

func setStringF(row, col int, s string, sh *xlsx.Sheet) func() error {
	return func() error {
		return setString(row, col, s, sh)
	}
}

func setCurrencyF(row, col int, d int, sh *xlsx.Sheet) func() error {
	return func() error {
		return setCurrency(row, col, d, sh)
	}
}
