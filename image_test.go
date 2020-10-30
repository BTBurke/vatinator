package vat

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStamp(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	var f bytes.Buffer
	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})

	if err := png.Encode(&f, img); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestHeader(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	var f bytes.Buffer

	img := createHeader(1, 10)

	if err := png.Encode(&f, img); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestComposite(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	f, err := os.Open("./test_receipts/PXL_20201002_163234312.jpg")
	require.NoError(t, err)
	defer f.Close()

	rcpt, err := jpeg.Decode(f)
	require.NoError(t, err)

	img := CompositeReceipt(1, rcpt, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)

	var out bytes.Buffer

	if err := png.Encode(&out, img); err != nil {
		assert.NoError(t, err)
	}
	snap.Assert(t, out.Bytes())
}

func TestPDF(t *testing.T) {
	ign := `@@\s\-3728,2\s\+3728,2\s@@`
	snap, err := snapshot.New(snapshot.IgnoreRegex(ign), snapshot.ContextLines(0))
	require.NoError(t, err)

	fname := "test.pdf"

	f, err := os.Open("./test_receipts/PXL_20201002_163234312.jpg")
	require.NoError(t, err)
	defer f.Close()

	rcpt, err := jpeg.Decode(f)
	require.NoError(t, err)

	f2, err := os.Open("./test_receipts/PXL_20201002_163306793.jpg")
	require.NoError(t, err)
	defer f2.Close()

	rcpt2, err := jpeg.Decode(f2)
	require.NoError(t, err)

	img := CompositeReceipt(1, rcpt, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)
	img2 := CompositeReceipt(2, rcpt2, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)

	pdf := NewPDF(fname)
	if err := pdf.WriteReceipt(img); err != nil {
		assert.NoError(t, err)
	}
	if err := pdf.WriteReceipt(img2); err != nil {
		assert.NoError(t, err)
	}
	if err := pdf.Save(); err != nil {
		assert.NoError(t, err)
	}

	pFile, err := ioutil.ReadFile(fname)
	require.NoError(t, err)
	require.NoError(t, os.Remove(fname))

	snap.Assert(t, pFile)

}
