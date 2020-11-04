package vat

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"path"

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

// output should be 100x100 png with a 10pixel red border around blue center
func TestCrop(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 120, 120))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(20, 20, 100, 100), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.Point{0, 0}, draw.Src)

	var f bytes.Buffer

	imgCropped := CropImage(img, 20, 20, 100, 100)

	if err := png.Encode(&f, imgCropped); err != nil {
		t.FailNow()
	}
	snapshot.Assert(t, f.Bytes())
}

func TestRotateCW(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	var f bytes.Buffer
	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})
	imgR := RotateCW(img)

	if err := png.Encode(&f, imgR); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestRotateCCW(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	var f bytes.Buffer
	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})
	imgR := RotateCCW(img)

	if err := png.Encode(&f, imgR); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestRotateByExif(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false))
	require.NoError(t, err)

	tt := []string{
		"test_receipts/exif-6.jpg",
		"test_receipts/exif-8.jpg",
	}

	for _, tc := range tt {
		t.Run(tc, func(t *testing.T) {
			wd, err := os.Getwd()
			require.NoError(t, err)

			f, err := os.Open(path.Join(wd, tc))
			require.NoError(t, err)

			img, err := RotateByExif(f)
			require.NoError(t, err)

			var out bytes.Buffer
			if err := png.Encode(&out, img); err != nil {
				t.FailNow()
			}
			snap.Assert(t, out.Bytes())
		})
	}
}
