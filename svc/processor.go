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
func NewSingleProcessor(db *badger.DB, accountID string, batchID string, keyPath string) Processor {
	return &singleProcessor{
		accountID: accountID,
		batchID:   batchID,
		db:        db,
		keyPath:   keyPath,
	}
}

type singleProcessor struct {
	accountID string
	batchID   string
	db        *badger.DB
	keyPath   string
}

func (s *singleProcessor) Add(name string, image img.Image) error {
	return process(s.db, s.accountID, s.batchID, name, image, s.keyPath, nil)
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
	hooks     *Hooks
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
	// Path to the Vision API key
	KeyPath string
	// Hooks to execute before/after processing the batch and receipts
	Hooks *Hooks
}

func NewParallelProcessor(db *badger.DB, accountID string, batchID string, opts *ParallelOptions) Processor {
	if opts == nil {
		opts = &ParallelOptions{
			ReprocessOnRulesChange: true,
			NumProcs:               20,
			KeyPath:                ".cfg/key.json",
		}
	}

	if opts.Hooks != nil && opts.Hooks.BeforeStart != nil {
		opts.Hooks.BeforeStart()
	}

	ch := make(chan parallelTask, opts.NumProcs+5)

	wg := &sync.WaitGroup{}
	for i := 0; i < opts.NumProcs; i++ {
		wg.Add(1)
		go func(ch chan parallelTask, db *badger.DB, accountID string, batchID string) {
			defer wg.Done()
			for task := range ch {
				if err := process(db, accountID, batchID, task.name, task.image, opts.KeyPath, opts.Hooks); err != nil {
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
		hooks:     opts.Hooks,
	}
}

func (p *parallelProcessor) Add(name string, image img.Image) error {
	p.ch <- parallelTask{name: name, image: image}
	return nil
}

func (p *parallelProcessor) Wait() error {
	close(p.ch)
	p.wg.Wait()
	if p.hooks != nil && p.hooks.AfterEnd != nil {
		p.hooks.AfterEnd()
	}
	return nil
}

// process image and save image and result to database
func process(db *badger.DB, accountID string, batchID string, name string, image img.Image, keyPath string, hooks *Hooks) error {
	if hooks != nil && hooks.BeforeEach != nil {
		// TODO: figure out how to do before each
	}
	result, err := ocr.ProcessImage(image, keyPath)
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
		Filename:          name,
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
	if result.Excise != nil {
		receipt.IsExcise = true
		receipt.ExciseType = result.Excise.Type
		receipt.ExciseAmount = result.Excise.Amount
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

	if hooks != nil && hooks.AfterEach != nil {
		if err := hooks.AfterEach(receipt); err != nil {
			return err
		}
	}

	return nil
}

var _ Processor = &singleProcessor{}
