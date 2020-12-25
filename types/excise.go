package types

import (
	"fmt"
	"math"
	"strconv"
)

const GasTaxRate float64 = 0.563

// Excise is an entry in the excise reimbursement form
type Excise struct {
	Type    string
	Content string
	Amount  string
	Tax     int
	Arve    string
	Date    string
}

// AsMap is called before exporting this receipt to the excise form.  If the tax is not explicitly set,
// it will be calculated automatically based on the current rate.
func (e *Excise) AsMap(i int) map[string]string {
	var tax int
	if e.Tax == 0 {
		tax = calculateTax(e.Amount, GasTaxRate)
		e.Tax = tax
	} else {
		tax = e.Tax
	}
	return map[string]string{
		makeKey("type", i):    maybe(e.Type),
		makeKey("content", i): e.Content,
		makeKey("amount", i):  maybe(e.Amount),
		makeKey("excise", i):  Currency(tax).String(),
		makeKey("arve", i):    fmt.Sprintf("%s / %s", maybe(e.Arve), maybe(e.Date)),
	}
}

// ExciseMetadata at the top of the excise form
type ExciseMetadata struct {
	Embassy string
	Name    string
	Bank    string
	Date    string
}

func calculateTax(amt string, rate float64) int {
	amtF, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		return 0
	}
	return int(math.Ceil(rate * amtF * 100))
}

func makeKey(f string, i int) string {
	return fmt.Sprintf("%s%d", f, i)
}

func maybe(s string) string {
	switch {
	case len(s) > 0:
		return s
	default:
		return "???"
	}
}
