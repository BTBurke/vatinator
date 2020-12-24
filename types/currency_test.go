package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	tt := []struct {
		in  int
		out string
	}{
		{0, "0.00"},
		{5, "0.05"},
		{10, "0.10"},
		{110, "1.10"},
		{15432, "154.32"},
	}
	for _, tc := range tt {
		assert.Equal(t, tc.out, Currency(tc.in).String())
	}
}
