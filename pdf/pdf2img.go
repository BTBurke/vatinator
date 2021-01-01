package pdf

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// PdfToImage converts a multipage pdf to a single image suitable for OCR like other
// receipts
// TODO: Possibly a problem if the PDF has zillions of pages and creates a huge image
func PdfToImage(r io.ReadCloser) (io.ReadCloser, error) {
	pdftoppmBin, err := exec.LookPath("pdftoppm")
	if err != nil {
		return nil, errors.Wrap(err, "requires pdftoppm to convert pdf to image")
	}
	convertBin, err := exec.LookPath("convert")
	if err != nil {
		return nil, errors.Wrap(err, "requires imagemagick to convert pdf to image")
	}

	// move input pdf to file in tempdir
	tmpdir, err := ioutil.TempDir("", "pdf2img")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tempdir for pdf2img")
	}
	pdfPath := filepath.Join(tmpdir, "in.pdf")
	targetPDF, err := os.Create(pdfPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pdf tmpfile in pdf2img")
	}
	if _, err := io.Copy(targetPDF, r); err != nil {
		return nil, errors.Wrap(err, "failed to copy input pdf to temp file in pdf2img")
	}

	// create series of images from pages of PDF
	outputPngs := filepath.Join(tmpdir, "out")
	pdfCmd := []string{"-png", "-rx", "300", "-ry", "300", pdfPath, outputPngs}
	cmd1 := exec.Command(pdftoppmBin, pdfCmd...)
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert pdf to images with output: %s", output1)
	}

	// convert multiple images to single image
	inputPngs := filepath.Join(tmpdir, "out*.png")
	outPng := filepath.Join(tmpdir, "result.png")
	convertCmd := []string{inputPngs, "-trim", "-append", outPng}
	cmd2 := exec.Command(convertBin, convertCmd...)
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert pdf images to single image with output: %s", output2)
	}

	f, err := os.Open(outPng)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get handle to the resulting image from pdf2img")
	}

	return f, nil
}
