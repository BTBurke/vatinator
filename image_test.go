package vat

import (
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStamp(t *testing.T) {
	f, err := os.Create("test.png")
	require.NoError(t, err)
	defer f.Close()

	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})

	if err := png.Encode(f, img); err != nil {
		t.FailNow()
	}

}

func TestHeader(t *testing.T) {
	f, err := os.Create("header.png")
	require.NoError(t, err)
	defer f.Close()

	img := createHeader(1, 10)

	if err := png.Encode(f, img); err != nil {
		t.FailNow()
	}

}

func TestComposite(t *testing.T) {
	f, err := os.Open("./test_receipts/PXL_20201002_163234312.jpg")
	require.NoError(t, err)

	//rcpt := image.NewRGBA(image.Rect(0, 0, 100, 100))
	//blue := color.RGBA{0, 0, 255, 255}
	//draw.Draw(rcpt, rcpt.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	defer f.Close()

	rcpt, err := jpeg.Decode(f)
	require.NoError(t, err)

	img := CompositeReceipt(1, rcpt, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 450)

	out, err := os.Create("composite.png")
	require.NoError(t, err)
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		assert.NoError(t, err)
	}
}

func TestPDF(t *testing.T) {
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

}
