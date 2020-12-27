package handlers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentStorageFile(t *testing.T) {
	data := []byte("test file")
	dataFile := fmt.Sprintf("%x", sha256.Sum256(data)) + ".jpg"

	dataPath, err := ioutil.TempDir("", "test-storage")
	require.NoError(t, err)
	defer os.RemoveAll(dataPath)

	r := ioutil.NopCloser(bytes.NewReader(data))
	if err := storeFileContent(r, dataPath, jpg); err != nil {
		assert.NoError(t, err)
	}
	finfo, err := os.Stat(filepath.Join(dataPath, dataFile))
	assert.NoError(t, err)

	// call again with same values should be no-op and fileinfo should remain the same
	r2 := ioutil.NopCloser(bytes.NewReader(data))
	if err := storeFileContent(r2, dataPath, jpg); err != nil {
		assert.NoError(t, err)
	}
	finfo2, err := os.Stat(filepath.Join(dataPath, dataFile))
	assert.NoError(t, err)
	assert.Equal(t, finfo, finfo2)
}

func TestContentStorageZip(t *testing.T) {
	// runs twice to test both single level and nested hierarchy zip files
	tt := []bool{false, true}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("nested %v", tc), func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "test-zip")
			require.NoError(t, err)
			defer os.RemoveAll(tmpdir)

			zipfile := filepath.Join(tmpdir, "test.zip")
			expected := createTestZip(zipfile, tc)

			datadir, err := ioutil.TempDir("", "test-zip-data")
			require.NoError(t, err)

			f, err := os.Open(zipfile)
			require.NoError(t, err)
			defer f.Close()

			assert.NoError(t, storeFileContent(f, datadir, zip))
			for _, e := range expected {
				path := filepath.Join(datadir, e)
				_, errF := os.Stat(path)
				assert.NoError(t, errF)
			}
		})
	}
}

// creates a test zip file with a bunch of fake jpgs with either 1 or 2 levels of directory hierarchy
func createTestZip(filename string, nested bool) []string {
	tmpdir, err := ioutil.TempDir("", "test-zip")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)

	// names of all the hashed files
	var out []string
	// original names
	var toArchive []string
	for i := 0; i < 10; i++ {
		data := []byte(fmt.Sprintf("test data %d", i))
		dataFile := fmt.Sprintf("%x", sha256.Sum256(data)) + ".jpg"

		var fname string
		switch nested {
		case true:
			// creates a nested hierarchy to test flattening
			fname = filepath.Join(tmpdir, fmt.Sprintf("dir%d", i), fmt.Sprintf("%d.jpg", i))
			_ = os.MkdirAll(filepath.Join(tmpdir, fmt.Sprintf("dir%d", i)), 0700)
		default:
			fname = filepath.Join(tmpdir, fmt.Sprintf("%d.jpg", i))
		}
		if err := ioutil.WriteFile(fname, data, 0644); err != nil {
			panic(err)
		}
		toArchive = append(toArchive, fname)
		out = append(out, dataFile)
	}
	if err := archiver.NewZip().Archive(toArchive, filename); err != nil {
		panic(err)
	}
	return out
}
