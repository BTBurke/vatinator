package svc

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"

	vat "github.com/BTBurke/vatinator"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/rs/xid"
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

// NewSingleProcessor returns a synchronous image processor that will run OCR and save the receipt
// results and the image to the database
func NewSingleProcessor(db *badger.DB, accountID string, batchID string) Processor {
	return &singleProcessor{
		accountID: accountID,
		batchID:   batchID,
		db:        db,
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

	//TODO: evaluate whether exif rotation is desired and what format to save
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
	var buf bytes.Buffer
	if err := png.Encode(&buf, croppedImage); err != nil {
		return errors.Wrapf(err, "failed to encode %s", name)
	}

	receipt := &Receipt{
		ID:            xid.New().String(),
		Vendor:        result.Vendor,
		TaxID:         result.TaxID,
		ReceiptNumber: result.ID,
		Total:         result.Total,
		VAT:           result.VAT,
		Date:          result.Date,
		BatchID:       batchID,
	}

	if err := db.Update(func(txn *badger.Txn) error {
		log.Printf("insert receipt")
		if err := upsertReceipt(txn, accountID, receipt); err != nil {
			return err
		}
		if err := upsertImage(txn, accountID, receipt.ID, buf.Bytes()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to persist receipt and image: %s", name)
	}

	return nil
}

func (s *singleProcessor) Wait() error {
	// returns immediately - synchronous
	return nil
}

var _ Processor = &singleProcessor{}
