package ocr

// Rule is an interface for all text processors that find a particular value from the raw
// text in a receipt
type Rule interface {
	Find(r *Result, text []string) error
}
