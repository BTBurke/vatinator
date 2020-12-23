package pdf

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/BTBurke/vatinator/img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPDF(t *testing.T) {
	ign := `@@\s\-6012,2\s\+6012,2\s@@`
	snap, err := snapshot.New(snapshot.IgnoreRegex(ign), snapshot.ContextLines(0))
	require.NoError(t, err)

	fname := "test.pdf"

	f, err := os.Open("../test_receipts/PXL_20201002_163234312.jpg")
	require.NoError(t, err)
	defer f.Close()

	rcpt, err := img.NewImageFromReader(f)
	require.NoError(t, err)

	f2, err := os.Open("../test_receipts/PXL_20201002_163306793.jpg")
	require.NoError(t, err)
	defer f2.Close()

	rcpt2, err := img.NewImageFromReader(f2)
	require.NoError(t, err)

	i, err := img.CompositeReceipt(1, rcpt, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)
	assert.NoError(t, err)
	i2, err := img.CompositeReceipt(2, rcpt2, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)
	assert.NoError(t, err)

	pdf := NewPDF(fname)
	if err := pdf.WriteReceipt(i); err != nil {
		assert.NoError(t, err)
	}
	if err := pdf.WriteReceipt(i2); err != nil {
		assert.NoError(t, err)
	}

	var b bytes.Buffer
	if err := pdf.Write(NopCloser(bufio.NewWriter(&b))); err != nil {
		assert.NoError(t, err)
	}

	snap.Assert(t, b.Bytes())

}

type nopCloser struct {
	io.Writer
}

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func (nopCloser) Close() error { return nil }
