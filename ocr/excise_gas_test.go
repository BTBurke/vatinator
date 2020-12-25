package ocr

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGasReceipt(t *testing.T) {
	tt := []struct {
		file   string
		amount string
		t      string
	}{
		{file: "testdata/gas1.txt", amount: "34.35", t: "95"},
		{file: "testdata/gas2.txt", amount: "39.80", t: "95"},
	}

	for _, tc := range tt {
		data, err := ioutil.ReadFile(tc.file)
		require.NoError(t, err)

		lines := strings.Split(string(data), "\n")
		isGas, amount, gasType := detectGasReceipt(lines)
		assert.True(t, isGas)
		assert.Equal(t, tc.amount, amount, fmt.Sprintf("file: %s", tc.file))
		assert.Equal(t, tc.t, gasType, fmt.Sprintf("file: %s", tc.file))
	}
}
