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
		{name: "bold food", input: "document no. 3188271b-eb74-44aa-8251- fa523af5242d", expect: "3188271b-eb74"},
		{name: "wolt", input: "order id: 601c1721af7e37fd4f032954", expect: "601c1721af7e37fd4f032954"},
		{name: "wolt2", input: "order id 601c1721af7e37fd4f032954", expect: "601c1721af7e37fd4f032954"},
		{name: "telia", input: "invoice 473892789234", expect: "473892789234"},
		{name: "alexela", input: "tsekk/arve 181638-5781", expect: "181638-5781"},
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
