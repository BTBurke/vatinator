package vatinator

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BTBurke/clt"
	"github.com/BTBurke/vatinator/bundled"
	"github.com/BTBurke/vatinator/img"
	"github.com/BTBurke/vatinator/svc"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
)

type Options struct {
	CredentialPath string
	OutputPath     string
	Interactive    bool
	log            *log.Logger
}

// ProcessService queues an async processing request for the web version.  CLI version calls
// Process directly.
type ProcessService interface {
	Do(id AccountID, batch string, month string, year int) error
}

type processService struct {
	exportDir string
	uploadDir string
	credFile  string
	account   AccountService
}

func NewProcessService(uploadDir string, exportDir string, credFile string, accountSvc AccountService) ProcessService {
	return processService{
		uploadDir: uploadDir,
		exportDir: exportDir,
		credFile:  credFile,
		account:   accountSvc,
	}
}

func (p processService) Do(id AccountID, batch string, month string, year int) error {
	path := filepath.Join(p.uploadDir, id.String(), batch)
	if finfo, err := os.Stat(path); err != nil || !finfo.IsDir() {
		return errors.Wrap(err, "could not find batch to process")
	}

	fdB, err := p.account.GetFormData(id)
	if err != nil {
		return errors.Wrapf(err, "could not get form data for account %s", id)
	}
	fd, err := UnmarshalFormData(fdB)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal form data for account %s", id)
	}

	opts := &Options{
		CredentialPath: p.credFile,
		OutputPath:     filepath.Join(path, "out"),
		Interactive:    false,
		log:            log.New(os.Stdout, batch, log.LstdFlags),
	}

	return Process(path, fd, month, year, opts)
}

// Process will read receipts located at path and process them into VAT and excise forms
func Process(path string, fd FormData, month string, year int, opts *Options) error {
	if opts == nil {
		opts = DefaultOptions(path)
	}
	if opts.log == nil {
		opts.log = log.New(ioutil.Discard, "", 0)
	}
	processStart := time.Now()
	opts.log.Printf("starting with config: %+v", opts)

	// best effort at clearing this directory, this can fail on windows so dont
	// check error.  Doesnt fuck up anything usually.
	_ = os.RemoveAll(opts.OutputPath)
	if err := os.MkdirAll(opts.OutputPath, 0700); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	template, err := bundled.Asset("assets/vat-template.xlsx")
	if err != nil {
		return errors.Wrap(err, "failed to find VAT form template")
	}

	// set up temporary database
	db, tempdir, err := createTempDB()
	if err != nil {
		return errors.Wrap(err, "failed to open temp database")
	}
	defer os.RemoveAll(tempdir)
	opts.log.Printf("created temporary database at %s", tempdir)

	// Walk directory and find all images to process
	var rcptFinder *clt.Progress
	if opts.Interactive {
		rcptFinder = clt.NewProgressSpinner("Finding receipts")
		rcptFinder.Start()
	}

	type task struct {
		name  string
		image img.Image
	}

	opts.log.Printf("looking for receipts in %s", path)
	var tasks []task
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		image, err := img.NewImageFromReader(f)
		if err != nil {
			return nil
		}
		tasks = append(tasks, task{path, image})

		return nil
	}); err != nil {
		if opts.Interactive {
			rcptFinder.Fail()
		}
		return errors.Wrap(err, "failed during scanning directory for images")
	}
	if opts.Interactive {
		rcptFinder.Success()
	}
	opts.log.Printf("found %d receipts in %s", len(tasks), path)

	months := map[string]int{"January": 1, "February": 2, "March": 3, "April": 4, "May": 5, "June": 6, "July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12}
	monthInt := months[month]

	// Start parallel processor and wait until finished

	var it *clt.Progress
	if opts.Interactive {
		it = clt.NewIncrementalProgressBar(len(tasks), "Doing magic to extract data from receipts")
		it.Start()

	}

	// these numbers dont matter, this is only a temp database
	accountID, batchID := "1", "1"

	errorWriter := svc.WriteErrors(filepath.Join(opts.OutputPath, "errors.txt"))
	proc := svc.NewParallelProcessor(db, accountID, batchID, &svc.ParallelOptions{
		ReprocessOnRulesChange: true,
		NumProcs:               20,
		Hooks: &svc.Hooks{
			AfterEach: func(r *svc.Receipt) error {
				if opts.Interactive {
					it.Increment()
				}
				if err := errorWriter(r); err != nil {
					return err
				}
				return nil
			},
		},
		KeyPath: opts.CredentialPath,
	})

	start := time.Now()
	for _, task := range tasks {
		if err := proc.Add(task.name, task.image); err != nil {
			return errors.Wrapf(err, "failed when processing %s", task.name)
		}
	}
	if err := proc.Wait(); err != nil {
		if opts.Interactive {
			it.Fail()
		}
		errors.Wrap(err, "processing images failed")
	}
	if opts.Interactive {
		it.Success()
	}
	opts.log.Printf("finished processing images in %s", time.Since(start))

	// export images to PDFs and fill forms
	var exp *clt.Progress
	if opts.Interactive {
		exp = clt.NewProgressSpinner("Filling out your forms")
		exp.Start()
	}

	opts.log.Printf("filling forms with data: %+v", fd)
	export := svc.NewExportService(db)
	if err := export.Create(accountID, batchID, &svc.ExportOptions{
		FirstName:    fd.FirstName,
		LastName:     fd.LastName,
		FullName:     fd.FullName,
		DiplomaticID: fd.DiplomaticID,
		Embassy:      fd.Embassy,
		Month:        month,
		MonthInt:     monthInt,
		Year:         year,
		Stamp:        []string{fd.FullName, fd.Embassy, fd.Address},
		Bank:         fd.Bank,
		Template:     template,
		OutputDir:    opts.OutputPath,
	}); err != nil {
		if opts.Interactive {
			exp.Fail()
		}
		return errors.Wrap(err, "failed during export")
	}
	if opts.Interactive {
		exp.Success()
	}
	opts.log.Printf("finished processing in %s", time.Since(processStart))

	return nil
}

func DefaultOptions(path string) *Options {
	// default options are set for CLI ops
	return &Options{
		CredentialPath: ".cfg/key.json",
		OutputPath:     filepath.Join(path, "out"),
		Interactive:    true,
		log:            log.New(os.Stdout, "", log.LstdFlags),
	}
}

func createTempDB() (*badger.DB, string, error) {
	tmpdir, err := ioutil.TempDir("", "vat")
	if err != nil {
		return nil, "", err
	}

	opts := badger.DefaultOptions(tmpdir)
	opts.Logger = nilLogger{}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, "", err
	}
	return db, tmpdir, nil
}

// gets rid of badger logging messages in terminal
type nilLogger struct{}

func (nilLogger) Errorf(s string, v ...interface{})   {}
func (nilLogger) Warningf(s string, v ...interface{}) {}
func (nilLogger) Infof(s string, v ...interface{})    {}
func (nilLogger) Debugf(s string, v ...interface{})   {}

var _ badger.Logger = nilLogger{}
