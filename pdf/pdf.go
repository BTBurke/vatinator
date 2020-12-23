package pdf

import (
	"fmt"
	"io"

	"github.com/BTBurke/vatinator/img"
	"github.com/jung-kurt/gofpdf"
)

// PDF handles creating PDFs
type PDF struct {
	Name        string
	p           *gofpdf.Fpdf
	numReceipts int
}

// NewPDF will create a new PDF handler with filename
func NewPDF(fname string) *PDF {
	p := &PDF{
		p:    gofpdf.New("P", "pt", "Letter", ""),
		Name: fname,
	}
	p.p.SetX(72.0)
	p.p.SetY(72.0)
	return p
}

// WriteReceipt will write a receipt to a new page in the PDF
func (p *PDF) WriteReceipt(image img.Image) error {

	r, err := image.NewReader()
	if err != nil {
		return err
	}

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
	itype.SetDpi(96.0)
	w, h := itype.Extent()

	p.p.ImageOptions(fmt.Sprintf("k%d", p.numReceipts+1), 72.0, 72.0, w, h, true, opt, 0, "")
	if !p.p.Ok() {
		return fmt.Errorf("Error while creating receipt pdf file: %s", p.p.Error())
	}
	p.numReceipts++

	return nil
}

// Save will write to the PDF file and close it
func (p *PDF) Save() error {
	return p.p.OutputFileAndClose(p.Name)
}

// Write will write the PDF file to the destination
func (p *PDF) Write(w io.WriteCloser) error {
	return p.p.OutputAndClose(w)
}
