package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExciseTaxCalculation(t *testing.T) {
	tt := []struct {
		in  string
		out int
	}{
		{in: "40.0", out: 2252},
		{in: "40", out: 2252},
		{in: "72.8", out: 4099},
	}
	for _, tc := range tt {
		assert.Equal(t, tc.out, calculateTax(tc.in, GasTaxRate))
	}
}
