package svc

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/BTBurke/vatinator/img"
	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageStorage(t *testing.T) {
	i := image.NewRGBA(image.Rect(0, 0, 120, 120))
	draw.Draw(i, i.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{0, 0}, draw.Src)

	testImage, err := img.NewImageFromImage(i)
	require.NoError(t, err)

	db, err := badger.Open(badger.DefaultOptions(t.TempDir()))
	require.NoError(t, err)

	imageService := NewImageService(db)
	if err := imageService.Upsert("test", "test", testImage); err != nil {
		assert.NoError(t, err)
	}

	recvImage, err := imageService.Get("test", "test")
	assert.NoError(t, err)
	assert.Equal(t, testImage, recvImage)

}
