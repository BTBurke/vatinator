package ocr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendor(t *testing.T) {
	tt := []struct {
		name string
		in   string
		out  string
	}{
		{name: "company at front", in: "OU APRANGA Estonia", out: "OU APRANGA Estonia"},
		{name: "with symbols", in: "H&M Hennes & Mauritz OÜ", out: "H&M Hennes & Mauritz OÜ"},
		{name: "misdetected 0 for O", in: "Test Company 0U", out: "Test Company OU"},
		{name: "COOP Ühistu", in: "COOP Ühistu", out: "COOP Ühistu"},
	}
	for _, tc := range tt {
		out := extractVendor([]string{tc.in})
		assert.Equal(t, tc.out, out)
	}
}
