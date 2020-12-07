package svc

import (
	"log"
	"os"
	"strings"
	"sync"

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
func (s *singleProcessor) Wait() error {
	// returns immediately - synchronous
	return nil
}

type parallelProcessor struct {
	wg        *sync.WaitGroup
	accountID string
	batchID   string
	db        *badger.DB
	numProcs  int
	ch        chan parallelTask
	reprocess bool
}

type parallelTask struct {
	name  string
	image img.Image
}

// ParallelOptions set options on a parallel image processor
type ParallelOptions struct {
	// If image has already been processed but rules have since changed, whether to reprocess using the new rules (default: true)
	ReprocessOnRulesChange bool
	// Number of images to process in parallel (default: 20)
	NumProcs int
}

func NewParallelProcessor(db *badger.DB, accountID string, batchID string, opts *ParallelOptions) Processor {
	if opts == nil {
		opts = &ParallelOptions{
			ReprocessOnRulesChange: true,
			NumProcs:               20,
		}
	}
	ch := make(chan parallelTask, opts.NumProcs+5)

	wg := &sync.WaitGroup{}
	for i := 0; i < opts.NumProcs; i++ {
		wg.Add(1)
		go func(ch chan parallelTask, db *badger.DB, accountID string, batchID string) {
			defer wg.Done()
			for task := range ch {
				if err := process(db, accountID, batchID, task.name, task.image); err != nil {
					log.Printf("processing error: %s", err)
				}
			}
		}(ch, db, accountID, batchID)
	}
	return &parallelProcessor{
		accountID: accountID,
		batchID:   batchID,
		db:        db,
		numProcs:  opts.NumProcs,
		ch:        ch,
		reprocess: opts.ReprocessOnRulesChange,
		wg:        wg,
	}
}

func (p *parallelProcessor) Add(name string, image img.Image) error {
	p.ch <- parallelTask{name: name, image: image}
	return nil
}

func (p *parallelProcessor) Wait() error {
	close(p.ch)
	p.wg.Wait()
	return nil
}

func process(db *badger.DB, accountID string, batchID string, name string, image img.Image) error {

	result, err := ocr.ProcessImage(image, "./vatinator-f91ccb107c2c.json")
	if err != nil {
		return errors.Wrapf(err, "failed to vision process %s", name)
	}

	f, err := os.Create(name + ".txt")
	if err == nil {
		_, _ = f.Write([]byte(strings.Join(result.Lines, "\n")))
		f.Close()
	}

	croppedImage, err := img.CropImage(image, int(result.Crop.Top), int(result.Crop.Left), int(result.Crop.Bottom), int(result.Crop.Right))
	if err != nil {
		return errors.Wrap(err, "failed to crop image")
	}

	receipt := &Receipt{
		ID:                xid.New().String(),
		Vendor:            result.Vendor,
		TaxID:             result.TaxID,
		ReceiptNumber:     result.ID,
		Total:             result.Total,
		VAT:               result.VAT,
		Date:              result.Date,
		BatchID:           batchID,
		Errors:            result.Errors,
		RulesVersion:      ocr.RulesVersion,
		CurrencyPrecision: Digit2,
	}

	if err := db.Update(func(txn *badger.Txn) error {

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

var _ Processor = &singleProcessor{}
