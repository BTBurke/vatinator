package update

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

type OS string

const (
	Windows OS = "Windows"
	Linux      = "Linux"
	MacOS      = "MacOS"
)

type updater struct {
	os          OS
	haveVersion string
	newVersion  string
	checksumID  int64
	assetID     int64
	exists      bool
}

type Updater interface {
	Check() error
	Update() error
	Exists() bool
	NewVersion() string
}

func NewUpdater(version string, os OS) Updater {
	return &updater{
		os:          os,
		haveVersion: version,
	}
}

func (u *updater) Check() error {
	// TODO: add timeouts
	ctx := context.Background()

	var err error
	u.exists, u.newVersion, u.assetID, u.checksumID, err = checkUpdate(ctx, u.haveVersion, u.os)
	if err != nil {
		return errors.Wrap(err, "failed to get update information")
	}
	return nil
}

func (u *updater) Update() error {
	// TODO: add timeouts and more sanity checks
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return NonFatalError{errors.Wrap(err, "failed to get current directory for update")}
	}

	tmpdir, err := ioutil.TempDir("", "vat_update")
	if err != nil {
		return NonFatalError{errors.Wrap(err, "failed to create temporary directory to download update")}
	}

	// find paths for assets
	checksumFile := filepath.Join(tmpdir, "checksum.txt")
	var archiveFile string
	var runningBinary string
	var replacementBinary string
	switch u.os {
	case Windows:
		archiveFile = filepath.Join(tmpdir, "update.zip")
		runningBinary = filepath.Join(wd, "vat.exe")
		replacementBinary = filepath.Join(tmpdir, "vat.exe")
	default:
		archiveFile = filepath.Join(tmpdir, "update.tar.gz")
		runningBinary = filepath.Join(wd, "vat")
		replacementBinary = filepath.Join(tmpdir, "vat")
	}

	// download checksums and update
	if err := downloadAsset(ctx, u.checksumID, checksumFile); err != nil {
		return NonFatalError{errors.Wrap(err, "failed to download checksum file for update")}
	}
	if err := downloadAsset(ctx, u.assetID, archiveFile); err != nil {
		return NonFatalError{errors.Wrap(err, "failed to download update")}
	}

	// check checksum for downloaded archive
	chkdata, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return NonFatalError{errors.Wrap(err, "failed to read checksum file")}
	}
	chk, err := getChecksum(chkdata, u.os)
	if err != nil {
		return NonFatalError{errors.Wrap(err, "failed to get checksum for this update")}
	}
	ok, err := checksum(archiveFile, chk)
	if err != nil {
		return NonFatalError{errors.Wrap(err, "error checking integrity of downloaded file")}
	}
	if !ok {
		return NonFatalError{fmt.Errorf("integrity check of update failed")}
	}

	// unarchive file and replace running binary
	if err := unarchive(archiveFile); err != nil {
		return NonFatalError{errors.Wrap(err, "failed to decompress the update")}
	}

	if _, err := replace(replacementBinary, runningBinary); err != nil {
		// errors here are already wrapped in Fatal or NonFatal, should pass them through
		// because it depends on where it happens in the update process
		return err
	}

	return nil
}

func (u *updater) Exists() bool {
	return u.exists
}

func (u *updater) NewVersion() string {
	return u.newVersion
}

func checkUpdate(ctx context.Context, haveVersion string, os OS) (exists bool, version string, assetID int64, checksumID int64, err error) {
	client := github.NewClient(nil)
	if ctx == nil {
		ctx = context.Background()
	}

	rel, _, err := client.Repositories.GetLatestRelease(ctx, "BTBurke", "vatinator")
	if err != nil {
		return
	}

	if haveVersion != rel.GetTagName() {
		for _, release := range rel.Assets {
			if strings.Contains(release.GetName(), string(os)) {
				exists = true
				version = rel.GetTagName()
				assetID = release.GetID()
			}
			if strings.Contains(release.GetName(), "checksums") {
				checksumID = release.GetID()
			}
		}
	}
	if version == "" || checksumID == 0 || assetID == 0 {
		return false, "", 0, 0, fmt.Errorf("missing a required file to update")
	}
	return
}

func downloadAsset(ctx context.Context, id int64, path string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	client := github.NewClient(nil)

	assetReader, _, err := client.Repositories.DownloadReleaseAsset(ctx, "BTBurke", "vatinator", id, http.DefaultClient)
	if err != nil || assetReader == nil {
		return err
	}
	defer assetReader.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, assetReader); err != nil {
		return err
	}

	return nil
}

func getChecksum(checksumFile []byte, os OS) (string, error) {
	checksums := strings.Split(string(checksumFile), "\n")
	for _, checksumLine := range checksums {
		if strings.Contains(checksumLine, string(os)) {
			return strings.Split(checksumLine, " ")[0], nil
		}
	}
	return "", fmt.Errorf("no checksum found for %s", os)
}

func unarchive(path string) error {
	dir, _ := filepath.Split(path)
	return archiver.Unarchive(path, dir)
}

func checksum(path string, checksum string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	got := fmt.Sprintf("%x", h.Sum(nil))
	return checksum == got, nil
}

func replace(src, dst string) (int64, error) {
	// check updated binary for sanity
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, NonFatalError{err}
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, NonFatalError{fmt.Errorf("%s is not a regular file", src)}
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, NonFatalError{err}
	}
	defer source.Close()

	// rename running binary to .old
	dstOld := dst + ".old"
	if err := os.Rename(dst, dstOld); err != nil {
		return 0, FatalError{err}
	}

	// create a new file with same name as running binary and copy updated binary
	destination, err := os.Create(dst)
	if err != nil {
		return 0, FatalError{err}
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if nBytes == 0 || err != nil {
		// updating failed, so copy old back.  If that fails, have to bail.  Everything is fucked.
		if err := os.Rename(dstOld, dst); err != nil {
			_ = os.Remove(dstOld)
			_ = os.Remove(dst)
			return 0, FatalError{fmt.Errorf("failed to undo a borked update")}
		}
	} else {
		// happy path, binary has been replaced, just have to remove the old one
		if err := os.Chmod(dst, 0777); err != nil {
			return 0, FatalError{err}
		}
		if err := os.Remove(dstOld); err != nil {
			return 0, FatalError{err}
		}
	}
	return nBytes, nil
}
