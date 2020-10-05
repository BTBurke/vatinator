package vat

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/jung-kurt/gofpdf"
)

// PDF handles creating PDFs
type PDF struct {
	p           *gofpdf.Fpdf
	fname       string
	numReceipts int
}

// NewPDF will create a new PDF handler with filename
func NewPDF(fname string) *PDF {
	p := &PDF{
		p:     gofpdf.New("P", "pt", "Letter", ""),
		fname: fname,
	}
	p.p.SetX(96.0)
	p.p.SetY(96.0)
	return p
}

// WriteReceipt will write a receipt to a new page in the PDF
func (p *PDF) WriteReceipt(receipt *image.RGBA) error {

	var buf bytes.Buffer
	if err := png.Encode(&buf, receipt); err != nil {
		return err
	}
	r := bytes.NewReader(buf.Bytes())

	opt := gofpdf.ImageOptions{
		ImageType:             "PNG",
		ReadDpi:               true,
		AllowNegativePosition: false,
	}

	p.p.AddPage()
	itype := p.p.RegisterImageOptionsReader(fmt.Sprintf("k%d", p.numReceipts+1), opt, r)
	if !p.p.Ok() {
		return fmt.Errorf("error while registering image: %s", p.p.Error())
	}
	w, h := itype.Extent()
	log.Printf("w: %f h: %f", w, h)

	p.p.ImageOptions(fmt.Sprintf("k%d", p.numReceipts+1), 72.0, 10.0, w, h, true, opt, 0, "")
	if !p.p.Ok() {
		return fmt.Errorf("Error while creating receipt pdf file: %s", p.p.Error())
	}
	p.numReceipts++

	return nil
}

// Save will write to the PDF file and close it
func (p *PDF) Save() error {
	return p.p.OutputFileAndClose(p.fname)
}
