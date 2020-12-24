package types

import "fmt"

// Excise is an entry in the excise reimbursement form
type Excise struct {
	Type    string
	Content string
	Amount  string
	Tax     int
	Arve    string
	Date    string
}

func (e Excise) AsMap(i int) map[string]string {
	return map[string]string{
		makeKey("type", i):    maybe(e.Type),
		makeKey("content", i): maybe(e.Content),
		makeKey("amount", i):  maybe(e.Amount),
		makeKey("excise", i):  Currency(e.Tax).String(),
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
