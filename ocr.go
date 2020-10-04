package vat

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var currency *regexp.Regexp
var kviitung *regexp.Regexp
var arve *regexp.Regexp
var tax *regexp.Regexp
var d *regexp.Regexp
var wellKnown map[string]string

const lineDither = int32(10)

func init() {
	currency = regexp.MustCompile(`[0-9]+\,[0-9]{2}`)
	kviitung = regexp.MustCompile(`kviitung[^0-9]+([0-9]*\/?[0-9]*)`)
	arve = regexp.MustCompile(`arve[^0-9]+([0-9]*)`)
	tax = regexp.MustCompile(`EE\s?[0-9]{9}`)
	d = regexp.MustCompile(`(01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31)\.?(01|02|03|04|05|06|07|08|09|10|11|12)\.?(2020|2021|20|21)`)
}

type Result struct {
	raw    []*pb.EntityAnnotation
	lines  []string
	File   string
	Date   string
	Total  int
	Tax    int
	Vendor string
	TaxID  string
	ID     string
	Crop   Crop
}

type Crop struct {
	Top    int32
	Bottom int32
	Right  int32
	Left   int32
}

func ProcessImage(fname string) (*Result, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening image %s: %v", fname, err)
	}
	img, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("error reading image %s: %v", fname, err)
	}
	ctx := context.Background()

	c, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("./vatinator-f91ccb107c2c.json"))
	if err != nil {
		return nil, fmt.Errorf("error creating vision client: %v", err)
	}

	res, err := c.DetectTexts(ctx, img, &pb.ImageContext{LanguageHints: []string{"ET"}}, 30)
	if len(res) == 0 || err != nil {
		return nil, fmt.Errorf("error detecting text: %v", err)
	}

	crop := getCrop(res)

	lines := strings.Split(res[0].Description, "\n")
	for _, l := range lines {
		lines = append(lines, strings.ToLower(l))
	}

	lines = joinFollowing(lines)
	extraLines := joinBigFuckingColumns(res)
	for _, l := range extraLines {
		extraLines = append(extraLines, strings.ToLower(l))
	}
	lines = append(lines, extraLines...)

	currencies := extractCurrency(lines)
	log.Printf("\ncurrencies: %v\n", currencies)
	tax, total, _ := extractTaxTotal(currencies)

	r := &Result{
		Crop:  crop,
		raw:   res,
		lines: lines,
		Tax:   tax,
		Total: total,
		TaxID: extractTaxID(lines),
		Date:  extractDate(lines),
		ID:    extractID(lines),
	}

	return r, nil

}

func getCrop(raw []*pb.EntityAnnotation) Crop {
	c := &Crop{
		Top:    math.MaxInt32,
		Bottom: int32(0),
		Left:   int32(0),
		Right:  math.MaxInt32,
	}

	for _, e := range raw {
		for _, v := range e.BoundingPoly.Vertices {
			if v.Y < c.Top {
				c.Top = v.Y
			}
			if v.Y > c.Bottom {
				c.Bottom = v.Y
			}
			if v.X < c.Right {
				c.Right = v.X
			}
			if v.X > c.Left {
				c.Left = v.X
			}
		}
	}
	return *c
}

// extracts all numbers of the form dd+,dd and returns them as integers in unit values (x100)
func extractCurrency(raw []string) []int {
	out := make([]int, 0)
	for _, line := range raw {
		c := currency.FindAllString(line, -1)
		for _, c1 := range c {
			cUnit := strings.Replace(c1, ",", "", -1)
			cAsInt, err := strconv.Atoi(cUnit)
			if err != nil {
				continue
			}
			out = append(out, cAsInt)
		}
	}
	return out
}

func extractID(lines []string) string {
	k := idFinder(kviitung, lines)
	if k == "" {
		return idFinder(arve, lines)
	}
	return k
}

func idFinder(r *regexp.Regexp, lines []string) string {
	for _, line := range lines {
		k := r.FindAllStringSubmatch(line, -1)
		if len(k) > 0 && len(k[0]) == 2 {
			if len(k[0][1]) > 0 {
				return k[0][1]
			}
		}
	}
	return ""
}

func extractTaxID(raw []string) string {
	for _, line := range raw {
		id := tax.FindString(line)
		if id != "" {
			return id
		}
	}
	return ""
}

func extractDate(raw []string) string {
	for _, line := range raw {
		r := d.FindAllStringSubmatch(line, -1)
		if len(r) > 0 && len(r[0]) == 4 {
			return fmt.Sprintf("%s/%s/%s", r[0][1], r[0][2], r[0][3])
		}
	}
	return ""
}

// math magic to determine tax and total by checking for 20% tax for every number on receipt
// only works because the values are sorted and it starts looking at the number most likely to be total
// TODO: doesn't handle the 9% or 10% tax brackets but fuck it
func extractTaxTotal(in []int) (tax int, total int, err error) {
	sort.Ints(in)
	for i := len(in)/2 - 1; i >= 0; i-- {
		opp := len(in) - 1 - i
		in[i], in[opp] = in[opp], in[i]
	}

	for _, i := range in {
		total = i
		expectedTax := total - int(float64(total)/1.20)
		log.Printf("Total: %d, Expected: %d\n", total, expectedTax)
		for _, j := range in {
			if j >= expectedTax-1 && j <= expectedTax+1 {
				tax = j
				return
			}
		}
	}
	return 0, 0, fmt.Errorf("no valid tax math found")
}

// joins successive lines to find small columns
func joinFollowing(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)

	for i := range in[0 : len(in)-1] {
		out = append(out, strings.Join(in[i:i+1], " "))
	}
	return out
}

type ColEntry struct {
	text string
	x    int32
	y    int32
}

type ColEntries []ColEntry

func (c ColEntries) Len() int      { return len(c) }
func (c ColEntries) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// this orders entries by their spatial position, first checking the y position to see if it
// is on a similar line +/- lineDither.  If they are on the same approximate line, order by x.
func (c ColEntries) Less(i, j int) bool {
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
	var cols ColEntries
	for _, entity := range in[1:] {
		e := ColEntry{
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
func recursivelyUnfuckColumn(curr string, y int32, agg []string, rest []ColEntry) []string {
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
