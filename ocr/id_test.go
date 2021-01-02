package ocr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceiptNumber(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		lines  []string
		expect string
	}{
		{name: "selver", input: "kviitung: 45065/90212", expect: "45065/90212"},
		{name: "bauhaus", lines: []string{"kv-arve", "086778"}, expect: "086778"},
		{name: "partnerkaart", input: "tšeki number: 118288", expect: "118288"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			switch {
			case len(tc.input) > 0:
				out := extractID([]string{tc.input})
				assert.Equal(t, tc.expect, out)
			default:
				out := extractID(tc.lines)
				assert.Equal(t, tc.expect, out)
			}
		})
	}
}
