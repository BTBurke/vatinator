package vat

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/require"
)

// TODO: make this more robust beyond just smoke test
func TestProcess(t *testing.T) {
	fname := "./test_receipts/PXL_20201002_163306793.jpg"
	f, err := os.Open(fname)
	require.NoError(t, err)

	res, err := ProcessImage(f)
	require.NoError(t, err)

	resB, err := json.Marshal(res)
	snapshot.Assert(t, resB)
}
