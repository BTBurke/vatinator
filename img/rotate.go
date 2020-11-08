package img

import (
	"fmt"
	"image"
	"log"

	"github.com/rwcarlsen/goexif/exif"
)

const (
	cw int = iota
	ccw
)

func RotateCW(img Image) (Image, error) {
	return rotateImage(img, cw)
}

func RotateCCW(img Image) (Image, error) {
	return rotateImage(img, ccw)
}

func rotateImage(orig Image, angle int) (Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, orig.Bounds().Max.Y, orig.Bounds().Max.X))
	for i := 0; i < orig.Bounds().Max.X; i++ {
		for j := 0; j < orig.Bounds().Max.Y; j++ {
			x2, y2 := mapPoint(i, j, orig.Bounds().Max.X, orig.Bounds().Max.Y, angle)
			img.Set(x2, y2, orig.At(i, j))
		}
	}
	return NewImageFromImage(img)
}

func mapPoint(x, y, width, height int, angle int) (x2 int, y2 int) {
	switch angle {
	case cw:
		x2 = height - y
		y2 = x
	default:
		x2 = y
		y2 = width - x
	}
	return
}

// RotateByExif reads embedded EXIF data and rotates the image to the correct aspect
func RotateByExif(i Image) (Image, error) {
	r, err := i.NewReader()
	if err != nil {
		return Image{}, err
	}
	e, _ := exif.Decode(r)

	oField, err := e.Get(exif.Orientation)
	if err != nil {
		return Image{}, err
	}

	orientation, err := oField.Int(0)
	if err != nil {
		return Image{}, err
	}
	log.Printf("orientation: %d", orientation)

	switch orientation {
	case 1:
		return i, nil
	case 8:
		return RotateCCW(i)
	case 3:
		i2, err := RotateCCW(i)
		if err != nil {
			return i2, err
		}
		return RotateCCW(i2)
	case 6:
		return RotateCW(i)
	default:
		return Image{}, fmt.Errorf("unknown exif rotation: %d", orientation)
	}
}
