package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BTBurke/vatinator/pdf"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

type fileType int

const (
	unknown fileType = iota
	jpg
	png
	zip
	pdfFile
)

func (ft fileType) String() string {
	switch ft {
	case jpg:
		return "jpg"
	case png:
		return "png"
	case zip:
		return "zip"
	case pdfFile:
		return "pdf"
	default:
		return ""
	}
}

func FileAddHandler(basepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		batchID := r.URL.Query().Get("batch_id")
		if len(batchID) == 0 {
			handleError(w, http.StatusBadRequest, fmt.Errorf("no batch id provided"))
			return
		}

		datapath := filepath.Join(basepath, batchID)
		if err := os.MkdirAll(datapath, 0700); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to create upload directory"))
			return
		}

		// switch on content-type to figure out most appropriate file name
		if err := r.ParseMultipartForm(400 * 1 << 20); err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "could not parse form data"))
			return
		}
		file, finfo, err := r.FormFile("file")
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "no file found in data"))
			return
		}

		ftype := typeFromContentType(finfo.Header.Get("Content-Type"))

		if err := storeFileContent(file, datapath, ftype); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error while writing uploaded image to storage"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func typeFromContentType(ct string) fileType {
	switch ct {
	case "image/jpeg":
		return jpg
	case "image/png":
		return png
	case "application/zip":
		return zip
	case "application/pdf":
		return pdfFile
	default:
		return unknown
	}
}

// storeFileContent will store the data at datapath based on SHA256 content hash.  If it already exists,
// it is a no-op.  This keeps exactly one copy of each uploaded file to prevent multiples per batch.
func storeFileContent(r io.ReadCloser, datapath string, ftype fileType) error {
	// zip files call storeFileContent recursively with the files in the zip
	if ftype == zip {
		return storeZipContent(r, datapath)
	}
	// pdf files call storeFileContent recursively after pdf is converted to an image
	if ftype == pdfFile {
		return storePDF2ImgContent(r, datapath)
	}

	tmpfile, err := ioutil.TempFile("", "fup-*."+ftype.String())
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpfile.Name())

	// multi write to temp file and to SHA256 hasher for content-addressable storage
	h := sha256.New()
	mw := io.MultiWriter(tmpfile, h)
	if _, err := io.Copy(mw, r); err != nil {
		return err
	}
	r.Close()
	tmpfile.Close()

	contentName := fmt.Sprintf("%x", h.Sum(nil))
	// check if file already exists by content hash
	targetPath := filepath.Join(datapath, fmt.Sprintf("%s.%s", contentName, ftype))
	_, errF := os.Stat(targetPath)
	switch {
	case errF != nil && !os.IsNotExist(errF):
		// unexpected error
		return err
	case errF == nil:
		// already exists, nothing to do
		return nil
	default:
		// doesnt exist, copy to datapath
		return copyFile(targetPath, tmpfile.Name())
	}
}

// storeZipContent opens the zip, walks it, then saves all the images
func storeZipContent(r io.ReadCloser, datapath string) error {
	tmpfile, err := ioutil.TempFile("", "fup-*.zip")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpfile.Name())

	// copy zip to temp file
	if _, err := io.Copy(tmpfile, r); err != nil {
		return err
	}
	tmpfile.Close()
	r.Close()

	// walk each zipped file and store file content.  Directory hierarchies are followed but flattened.
	return archiver.NewZip().Walk(tmpfile.Name(), func(f archiver.File) error {
		if f.FileInfo.IsDir() {
			return nil
		}
		switch {
		case strings.HasSuffix(f.FileInfo.Name(), "jpg"):
			return storeFileContent(f, datapath, jpg)
		case strings.HasSuffix(f.FileInfo.Name(), "png"):
			return storeFileContent(f, datapath, png)
		default:
			return fmt.Errorf("unknown filetype: %s", f.FileInfo.Name())
		}
	})
}

// storePDF2ImgContent converts an input pdf to an image before storing it
func storePDF2ImgContent(r io.ReadCloser, datapath string) error {
	rc, err := pdf.PdfToImage(r)
	if err != nil {
		return errors.Wrap(err, "failed to convert uploaded pdf to image")
	}

	return storeFileContent(rc, datapath, png)
}

func copyFile(to, from string) error {
	target, err := os.Create(to)
	if err != nil {
		return err
	}
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	if _, err := io.Copy(target, source); err != nil {
		return err
	}
	return nil
}

func handleError(w http.ResponseWriter, status int, err error) {
	log.Printf("error: %s", err)
	w.WriteHeader(status)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.Printf("meta error, error writing error: %v", err)
	}
}
