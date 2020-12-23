package ocr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	tt := []struct {
		name  string
		in    []string
		tax   int
		total int
	}{
		{name: "mixed like selver", in: []string{"0,988", "0,956", "68,66", "13,73", "82,39", "69,26", "13,73", "82,99"}, tax: 1373, total: 8239},
		{name: "tax in currency3", in: []string{"4,00", "1,00", "2,00", "5,000", "1,000", "0,000", "6,00"}, tax: 100, total: 600},
		{name: "with spaces like H&M", in: []string{"16, 67", "83, 27", "99, 94"}, tax: 1667, total: 9994},
	}
	for _, tc := range tt {
		tax, total, _ := findTaxTotal(tc.in)
		assert.Equal(t, tc.tax, tax)
		assert.Equal(t, tc.total, total)
	}
}
