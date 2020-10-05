package vat

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

const HeaderHeight int = 20
const HeaderText string = "Kviitung %d"

var StampColor color.RGBA = color.RGBA{125, 3, 8, 125}

// createStamp creates a minimum size stamp with the text given in lines with a transparent background
// and text color given by StampColor.  If lines is empty, it will return a zero pixel image.
func createStamp(lines []string) *image.RGBA {
	if len(lines) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	// number of pixels on top and bottom to account for characters that
	// hang below the baseline
	topBtmPad := 5

	// lineheight set slightly smaller than text size based on trial and error
	lineHeight := 14

	max := 0
	for _, line := range lines {
		if len(line) > max {
			max = len(line)
		}
	}

	w := max * 8
	h := (len(lines) * lineHeight) + 2*topBtmPad

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(StampColor),
		Face: inconsolata.Bold8x16,
	}

	for i, line := range lines {
		point := fixed.Point26_6{Y: fixed.Int26_6(((lineHeight * (i + 1)) + topBtmPad) * 64)}
		d.Dot = point
		d.DrawString(line)
	}
	return img
}

// createHeader creates a header above the receipt using HeaderText as the format
// string.  The width, w, determines the maximum width, but it will grow automatically to
// fit the entire header text plus 8 pixels if it is too small.
func createHeader(num int, w int) *image.RGBA {
	text := fmt.Sprintf(HeaderText, num)
	textWidth := len(text) * 8

	// if label wider than receipt, then make header minimum width
	leftPad := (w - textWidth) / 2
	if leftPad < 0 {
		w = textWidth + 8
		leftPad = 4
	}

	img := image.NewRGBA(image.Rect(0, 0, w, HeaderHeight))
	for i := 0; i < w; i++ {
		for j := 0; j < HeaderHeight; j++ {
			img.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: inconsolata.Bold8x16,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(leftPad * 64), Y: fixed.Int26_6(16 * 64)},
	}
	d.DrawString(text)

	return img
}

// CompositeReceipt will create a composite image of the receipt appropriately scaled to fit on the page,
// with a header showing the receipt number, and the superimposed stamp.  If the stamp is empty or nil, no stamp
// is applied.  If stampY is 0, it will be placed in a default position about 1/3 down the receipt.
func CompositeReceipt(num int, receipt image.Image, stamp []string, stampY int) *image.RGBA {

	// Check receipt and resize to make it fit on letter size paper at 72 dpi
	rcptWidth := receipt.Bounds().Max.X
	rcptHeight := receipt.Bounds().Max.Y
	if rcptHeight > 648 {
		receipt = resize.Resize(0, 96*9, receipt, resize.Lanczos3)
		rcptWidth = receipt.Bounds().Max.X
		rcptHeight = receipt.Bounds().Max.Y
	}
	if rcptWidth > 468 {
		receipt = resize.Resize(0, 96*6.5, receipt, resize.Lanczos3)
		rcptWidth = receipt.Bounds().Max.X
		rcptHeight = receipt.Bounds().Max.Y
	}

	stampImg := createStamp(stamp)
	stampWidth := stampImg.Bounds().Max.X
	stampHeight := stampImg.Bounds().Max.Y

	headerImg := createHeader(num, rcptWidth)
	headerWidth := headerImg.Bounds().Max.X
	headerHeight := headerImg.Bounds().Max.Y

	// TODO: hacky way of getting max width, should fix this
	finalWidth := rcptWidth
	if headerWidth > finalWidth {
		finalWidth = headerWidth
	}
	if stampWidth > finalWidth {
		finalWidth = stampWidth
	}
	finalHeight := rcptHeight + headerHeight

	// padding for stamp
	if stampY == 0 {
		// default place 1/3 down the page
		stampY = rcptHeight / 3
	}
	stampLeft := (finalWidth - stampWidth) / 2
	stampTop := headerHeight + stampY - stampHeight
	if stampTop < 0 {
		stampTop = 0
	}

	img := image.NewRGBA(image.Rect(0, 0, finalWidth, finalHeight))
	draw.Draw(img, image.Rect(0, 0, headerWidth, headerHeight), headerImg, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(0, headerHeight, finalWidth, finalHeight), receipt, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(stampLeft, stampTop, stampLeft+stampWidth, stampTop+stampHeight), stampImg, image.Point{0, 0}, draw.Over)

	return img
}

// CropImage returns a copy of the image cropped to the top, right, bottom, and left point
func CropImage(orig *image.RGBA, top, right, bottom, left int) image.Image {
	return orig.SubImage(image.Rect(left, top, right, bottom))
}
