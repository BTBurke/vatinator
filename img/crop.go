package img

import (
	"image"
	"image/draw"
)

const CropPadding int = 10

// CropImage returns a copy of the image cropped to the top, right, bottom, and left point.  CropPadding setting
// will retain that number of pixels as a margin on all sides around the passed crop points.
func CropImage(orig Image, top, left, bottom, right int) (Image, error) {
	left = max(0, left-CropPadding)
	right = min(orig.Bounds().Max.X, right+CropPadding)
	top = max(0, top-CropPadding)
	bottom = min(orig.Bounds().Max.Y, bottom+CropPadding)

	out := image.NewRGBA(image.Rect(0, 0, (right - left), (bottom - top)))
	draw.Draw(out, out.Bounds(), orig, image.Point{left, top}, draw.Src)
	return NewImageFromImage(out)
}

func max(x, y int) int {
	switch {
	case x > y:
		return x
	default:
		return y
	}
}

func min(x, y int) int {
	switch {
	case x < y:
		return x
	default:
		return y
	}

}
