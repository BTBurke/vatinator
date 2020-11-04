package svc

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"

	vat "github.com/BTBurke/vatinator"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
)

var cache map[string]Processor

type Processor interface {
	Add(name string, r io.Reader) error
	Wait() error
}

func init() {
	if cache == nil {
		cache = make(map[string]Processor)
	}
}

type singleProcessor struct {
	accountID string
	batchID   string
	db        *badger.DB
}

func (s *singleProcessor) Add(name string, r io.Reader) error {
	return process(s.db, s.accountID, s.batchID, name, r)
}

func process(db *badger.DB, accountID string, batchID string, name string, r io.Reader) error {
	i, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read: %s", name)
	}

	//TODO: rotate based on exif, reencode as png
	rotatedImage, err := vat.RotateByExif(bytes.NewReader(i))
	if err != nil {
		return errors.Wrapf(err, "failed to rotate: %s", name)
	}
	var b bytes.Buffer
	if err := png.Encode(&b, rotatedImage); err != nil {
		return errors.Wrapf(err, "failed to encode %s", name)
	}

	// process with google cloud vision
	visionReader := bytes.NewReader(b.Bytes())

	result, err := vat.ProcessImage(visionReader)
	if err != nil {
		return errors.Wrapf(err, "failed to vision process %s", name)
	}

	// decode and crop image based on OCRed positions
	img, _, err := image.Decode(bytes.NewReader(i))
	if err != nil {
		return errors.Wrapf(err, "failed to decode %s", name)
	}
	croppedImage := vat.CropImage(img, int(result.Crop.Top), int(result.Crop.Left), int(result.Crop.Bottom), int(result.Crop.Right))

	// store image and results
	fmt.Println(croppedImage)

	return nil
}
