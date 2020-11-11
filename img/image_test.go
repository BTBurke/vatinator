package img

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStamp(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	var f bytes.Buffer
	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})

	if err := png.Encode(&f, img); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestHeader(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	var f bytes.Buffer

	img := createHeader(1, 10)

	if err := png.Encode(&f, img); err != nil {
		t.FailNow()
	}
	snap.Assert(t, f.Bytes())

}

func TestComposite(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	f, err := os.Open("../test_receipts/PXL_20201002_163234312.jpg")
	require.NoError(t, err)
	defer f.Close()

	rcpt, err := NewImageFromReader(f)
	require.NoError(t, err)

	img, err := CompositeReceipt(1, rcpt, []string{"Bryan Burke", "US Embassy", "Kentmanni 20"}, 0)
	assert.NoError(t, err)

	b, err := img.AsPNG()
	require.NoError(t, err)

	snap.Assert(t, b)
}

// output should be 100x100 png with a 10pixel red border around blue center
func TestCrop(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	img := image.NewRGBA(image.Rect(0, 0, 120, 120))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(20, 20, 100, 100), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.Point{0, 0}, draw.Src)
	i, err := NewImageFromImage(img)
	require.NoError(t, err)

	imgCropped, err := CropImage(i, 20, 20, 100, 100)
	require.NoError(t, err)

	b, err := imgCropped.AsPNG()
	require.NoError(t, err)

	snap.Assert(t, b)
}

func TestRotateCW(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})
	i, err := NewImageFromImage(img)
	require.NoError(t, err)

	imgR, err := RotateCW(i)
	assert.NoError(t, err)

	b, err := imgR.AsPNG()
	assert.NoError(t, err)

	snap.Assert(t, b)

}

func TestRotateCCW(t *testing.T) {
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	img := createStamp([]string{"test 1", "test 2", "fucking jelly 3"})
	i, err := NewImageFromImage(img)
	require.NoError(t, err)

	imgR, err := RotateCCW(i)
	require.NoError(t, err)

	b, err := imgR.AsPNG()
	require.NoError(t, err)

	snap.Assert(t, b)

}

// func TestRotateByExif(t *testing.T) {
// snap, err := snapshot.New(snapshot.Diffable(false))
// require.NoError(t, err)
//
// tt := []string{
// "../test_receipts/exif-6.jpg",
// "../test_receipts/exif-8.jpg",
// }
//
// for _, tc := range tt {
// t.Run(tc, func(t *testing.T) {
// wd, err := os.Getwd()
// require.NoError(t, err)
//
// f, err := os.Open(path.Join(wd, tc))
// require.NoError(t, err)
// i, err := NewImageFromReader(f)
// require.NoError(t, err)
//
// img, err := RotateByExif(i)
// require.NoError(t, err)
//
// b, err := img.AsPNG()
// require.NoError(t, err)
// snap.Assert(t, b)
// })
// }
// }
