package svc

import (
	"log"

	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/ocr"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var cache map[string]Processor

type Processor interface {
	Add(name string, image img.Image) error
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

func (s *singleProcessor) Add(name string, image img.Image) error {
	return process(s.db, s.accountID, s.batchID, name, image)
}

func process(db *badger.DB, accountID string, batchID string, name string, image img.Image) error {

	//TODO: evaluate whether exif rotation is desired and what format to save
	rotatedImage, err := img.RotateByExif(image)
	if err != nil {
		return errors.Wrapf(err, "failed to rotate: %s", name)
	}

	result, err := ocr.ProcessImage(rotatedImage)
	if err != nil {
		return errors.Wrapf(err, "failed to vision process %s", name)
	}

	croppedImage, err := img.CropImage(rotatedImage, int(result.Crop.Top), int(result.Crop.Left), int(result.Crop.Bottom), int(result.Crop.Right))
	if err != nil {
		return errors.Wrap(err, "failed to crop image")
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
		if err := upsertImage(txn, accountID, receipt.ID, croppedImage); err != nil {
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
