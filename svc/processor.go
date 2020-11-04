package svc

import (
	vat "github.com/BTBurke/vatinator"
	"sync"
)

var cache map[string]Processor

type Processor interface {
	Add(r io.Reader) error
	Wait() error
}

func init() {
	if cache == nil {
		cache = make(map[string]Processor)
	}
}

type singleProcessor struct {
	db *badger.DB
}

func (s *singleProcessor) Add(r io.Reader) error {
	return process(r)
}

func process(r io.Reader) error {
	image, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	visionReader := bytes.NewReader(image)

	result, err := vat.ProcessImage(visionReader)
	if err != nil {
		return err
	}

	// decode and crop image based on OCRed image
	img, err := image.Decode(image)
	if err != nil {
		return err
	}
	croppedImage, err := vat.CropImage(img, result.Crop.Top, result.Crop.Left, result.Crop.Bottom, result.Crop.Right)
	if err != nil {
		return err
	}
	fmt.Println(croppedImage)

	return nil
}
