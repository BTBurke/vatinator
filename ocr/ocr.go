package ocr

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/BTBurke/vatinator/img"
	"github.com/disintegration/imaging"

	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// RulesVersion is a magic string that is recorded with each receipt to denote which version of the
// extraction rules was used.  Reprocessing is the default when a receipt was processed under old rules.
// Format for RulesVersion is YYYYMMDD.  It doesn't matter what it is, so a version can be added for multiple changes
// on the same day (e.g., YYYYMMDD-v1).
// TODO: shift to some build time hash that denotes if the rules have changed
var RulesVersion string = "20201206"

// lineDither is the number of pixels in the Y direction that two words should be considered to be on the same
// line.  This is used to reconstruct multi-column receipt formats separated by large white space.
const lineDither = int32(10)

type CurrencyPrecision int

const (
	Currency2 CurrencyPrecision = iota + 2
)

type ReceiptType int

const (
	Regular = iota
	Gas
)

// Result
type Result struct {
	raw         []*pb.EntityAnnotation
	Lines       []string
	Orientation Orientation
	File        string
	// date format dd/mm/yy or dd/mm/yyyy depending on how it is detected on the receipt
	Date   string
	Total  int
	VAT    int
	Vendor string
	TaxID  string
	ID     string
	Excise *Excise
	Crop   Crop
	Errors []string
}

// Crop returns the pixel location of the tightest crop that contains all
// recognized text
type Crop struct {
	Top    int32
	Bottom int32
	Right  int32
	Left   int32
}

type Excise struct {
	Type   string
	Amount string
}

// ProcessImage uses a pre-trained ML model to extract text from the receipt image, then
// a series of regular expressions and text manipulation to find the VAT data
func ProcessImage(image img.Image, credPath string) (*Result, error) {

	res, err := doAnnotation(image, credPath)
	if err != nil {
		return nil, err
	}
	orient := DetectOrientation(res)
	if orient != Orientation0 {
		// redo detection after rotations so crop is right and I dont have to figure it out
		rotatedImage, err := AutoRotateImage(image, orient)
		if err != nil {
			return nil, err
		}
		res, err = doAnnotation(rotatedImage, credPath)
		if err != nil {
			return nil, err
		}
	}

	// find the minimum bounding box for the receipt
	crop := getCrop(res)

	// get normal case lines, then add lowercase version for each value
	lines := strings.Split(res[0].Description, "\n")
	for _, l := range lines {
		lines = append(lines, strings.ToLower(l))
	}

	// joins blocks of words that are close together to look for successive matches
	lines = joinFollowing(lines)

	// looks for wide columns by spatial comparison
	extraLines := joinBigFuckingColumns(res)
	for _, l := range extraLines {
		extraLines = append(extraLines, strings.ToLower(l))
	}
	lines = append(lines, extraLines...)

	// looks for little tables with a header and value
	extraLines2 := joinLittleTables(res)
	for _, l := range extraLines2 {
		extraLines2 = append(extraLines2, strings.ToLower(l))
	}
	lines = append(lines, extraLines2...)

	rules := []Rule{
		VendorRule(),
		DateRule(),
		IDRule(),
		CurrencyRule(),
		GasRule(),
	}

	r := &Result{
		Crop:        crop,
		Lines:       lines,
		Orientation: orient,
	}

	for _, rule := range rules {
		if err := rule.Find(r, lines); err != nil {
			return nil, err
		}
	}

	return r, nil

}

func AutoRotateImage(image img.Image, orient Orientation) (img.Image, error) {
	switch orient {
	case Orientation90:
		i, err := img.NewImageFromImage(imaging.Rotate90(image.GetImage()))
		if err != nil {
			return img.Image{}, err
		}
		return i, nil
	case Orientation180:
		i, err := img.NewImageFromImage(imaging.Rotate180(image.GetImage()))
		if err != nil {
			return img.Image{}, err
		}
		return i, nil
	case Orientation270:
		i, err := img.NewImageFromImage(imaging.Rotate270(image.GetImage()))
		if err != nil {
			return img.Image{}, err
		}
		return i, nil
	default:
		return img.Image{}, fmt.Errorf("unknown rotation")
	}
}

func doAnnotation(image img.Image, credPath string) ([]*pb.EntityAnnotation, error) {
	imgReader, err := image.NewReader()
	if err != nil {
		return nil, err
	}

	i, err := vision.NewImageFromReader(imgReader)
	if err != nil {
		return nil, fmt.Errorf("error reading image: %v", err)
	}
	ctx := context.Background()

	c, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, fmt.Errorf("error creating vision client: %v", err)
	}

	res, err := c.DetectTexts(ctx, i, &pb.ImageContext{LanguageHints: []string{"ET"}}, 1000)
	if len(res) == 0 || err != nil {
		return nil, fmt.Errorf("error detecting text: %v", err)
	}
	return res, nil
}

// determines the minimum bounding box for the text on the receipt
func getCrop(raw []*pb.EntityAnnotation) Crop {
	c := Crop{
		Top:    math.MaxInt32,
		Bottom: int32(0),
		Right:  int32(0),
		Left:   math.MaxInt32,
	}

	for _, e := range raw {
		for _, v := range e.BoundingPoly.Vertices {
			if v.Y < c.Top {
				c.Top = v.Y
			}
			if v.Y > c.Bottom {
				c.Bottom = v.Y
			}
			if v.X < c.Left {
				c.Left = v.X
			}
			if v.X > c.Right {
				c.Right = v.X
			}
		}
	}
	return c
}

// joins successive lines to find small columns
func joinFollowing(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)

	for i := range in[0 : len(in)-1] {
		out = append(out, strings.Join(in[i:i+2], " "))
	}
	return out
}

type colEntry struct {
	text string
	x    int32
	y    int32
}

// data structure to hold column entries.  Each potential column entry is compared based on X,Y coordinates
// to selective join them into a line.
type colEntries []colEntry

func (c colEntries) Len() int      { return len(c) }
func (c colEntries) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// this orders entries by their spatial position, first checking the y position to see if it
// is on a similar line +/- lineDither.  If they are on the same approximate line, order by x.
func (c colEntries) Less(i, j int) bool {
	switch {
	case c[i].y <= c[j].y-lineDither:
		return true
	case c[i].y >= c[j].y+lineDither:
		return false
	default:
		return c[i].x <= c[j].x
	}
}

// this looks for very wide columns that would normally be impossible to unfuck and makes them only
// nearly impossible to unfuck
func joinBigFuckingColumns(in []*pb.EntityAnnotation) []string {
	if len(in) <= 1 {
		return nil
	}
	var cols colEntries
	for _, entity := range in[1:] {
		e := colEntry{
			text: entity.Description,
			x:    entity.BoundingPoly.Vertices[0].X,
			y:    entity.BoundingPoly.Vertices[0].Y,
		}
		cols = append(cols, e)
	}
	sort.Sort(cols)
	var out []string
	return recursivelyUnfuckColumn(cols[0].text, cols[0].y, out, cols[1:])

}

// get out your compsci textbook kids, we are going recursive on this motherfucker to join
// extracted text that appears to be on the same line
func recursivelyUnfuckColumn(curr string, y int32, agg []string, rest []colEntry) []string {
	if len(rest) == 0 {
		return append(agg, curr)
	}

	// if its on the same line add it to the current and keep going
	if rest[0].y <= y+lineDither && rest[0].y >= y-lineDither {
		return recursivelyUnfuckColumn(fmt.Sprintf("%s %s", curr, rest[0].text), rest[0].y, agg, rest[1:])
	}
	// otherwise, add this line to the output, and start on a new line
	return recursivelyUnfuckColumn(rest[0].text, rest[0].y, append(agg, curr), rest[1:])
}

// looks for small tables with a header followed by a value like arve then value under
type tableEntry struct {
	text  string
	x     int32
	y     int32
	width int32
}

type tableEntries []tableEntry

func (t tableEntries) Len() int      { return len(t) }
func (t tableEntries) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// orders by spatial position, looking to see if two values substantially overlap
// in the x direction, indicating they should be in a column then ordering by their spatial
// position north to south
func (t tableEntries) Less(i, j int) bool {
	switch {
	case t[i].x <= t[j].x-t[i].width/2:
		return true
	case t[i].x >= t[j].x+t[i].width/2:
		return false
	default:
		return t[i].y <= t[j].y
	}
}

func joinLittleTables(in []*pb.EntityAnnotation) []string {
	if len(in) <= 1 {
		return nil
	}
	var tables tableEntries
	for _, entity := range in[1:] {
		e := tableEntry{
			text:  entity.Description,
			x:     entity.BoundingPoly.Vertices[0].X,
			y:     entity.BoundingPoly.Vertices[0].Y,
			width: entity.BoundingPoly.Vertices[1].X - entity.BoundingPoly.Vertices[0].X,
		}
		tables = append(tables, e)
	}
	sort.Sort(tables)
	var out []string
	return recursivelyFindTables(tables[0].text, tables[0].x, tables[0].width, out, tables[1:])
}

// look for things that are in the same column that are in similar spatial position to represent little tables
func recursivelyFindTables(curr string, x int32, w int32, agg []string, rest []tableEntry) []string {
	if len(rest) == 0 {
		return append(agg, curr)
	}

	// if it's in the same column, specified by +/- half the width of the current value, add it to the current line and keep going
	if rest[0].x <= x+w/2 && rest[0].x <= x-w/2 {
		return recursivelyFindTables(fmt.Sprintf("%s %s", curr, rest[0].text), rest[0].x, rest[0].width, agg, rest[1:])
	}

	// otherwise, add this line to output and start on a new line
	return recursivelyFindTables(rest[0].text, rest[0].x, rest[0].width, append(agg, curr), rest[1:])
}
