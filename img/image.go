package img

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"time"

	"github.com/BTBurke/vatinator/db"
)

type Image struct {
	b      []byte
	format string
	image  image.Image
}

func NewImageFromBytes(b []byte) (Image, error) {
	img, f, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return Image{}, err
	}

	return Image{
		b:      b,
		format: f,
		image:  img,
	}, nil
}

func NewImageFromImage(i image.Image) (Image, error) {
	b, err := encodePNG(i)
	if err != nil {
		return Image{}, err
	}

	return Image{
		b:      b,
		format: "png",
		image:  i,
	}, nil
}

func NewImageFromReader(r io.Reader) (Image, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return Image{}, err
	}
	return NewImageFromBytes(b)
}

func encodePNG(i image.Image) ([]byte, error) {
	var b bytes.Buffer
	if err := png.Encode(&b, i); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodeJPG(i image.Image) ([]byte, error) {
	var b bytes.Buffer
	if err := jpeg.Encode(&b, i, nil); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (i Image) AsPNG() ([]byte, error) {
	if i.b != nil && i.format == "png" {
		return i.b, nil
	}
	return encodePNG(i.image)
}

func (i Image) AsJPG() ([]byte, error) {
	if i.b != nil && i.format == "jpeg" {
		return i.b, nil
	}
	return encodeJPG(i.image)
}

func (i Image) NewReader() (io.Reader, error) {
	if len(i.format) > 0 && i.b != nil {
		return bytes.NewReader(i.b), nil
	}
	b, err := encodePNG(i.image)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Implements image.Image interface by wrapping underlying image

func (i Image) ColorModel() color.Model {
	return i.image.ColorModel()
}

func (i Image) Bounds() image.Rectangle {
	return i.image.Bounds()
}

func (i Image) At(x, y int) color.Color {
	return i.image.At(x, y)
}

// Implements db.Entity to marshal/unmarshal itself from embedded db

var DefaultImageDuration time.Duration = 24 * time.Hour * 95

func (i *Image) TTL() time.Duration {
	return DefaultImageDuration
}

func (i *Image) Type() byte {
	return db.Image
}

func (i *Image) MarshalBinary() ([]byte, error) {
	return i.AsPNG()
}

func (i *Image) UnmarshalBinary(data []byte) error {
	img, err := NewImageFromBytes(data)
	if err != nil {
		return err
	}
	*i = img
	return nil
}

var _ image.Image = Image{}
var _ db.Entity = &Image{}
