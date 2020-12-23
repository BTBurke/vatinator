package update

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const checksum_id_v16 int64 = 29907733
const checksum_v16_windows string = "5bb9398a96fd1d0cecf8c75a9e663dec07d5327d1fd3b95c8359742aaec7ff69"
const binary_windows_id_v16 int64 = 29907735

func TestCheck(t *testing.T) {
	exists, version, id, checksum, err := checkUpdate(context.Background(), "v1.0.0", Windows)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, id, int64(0))
	assert.GreaterOrEqual(t, checksum, int64(0))
	assert.GreaterOrEqual(t, len(version), 0)
}

func TestDownload(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "vat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	p := filepath.Join(tmpdir, "checksum.txt")

	if err := downloadAsset(nil, checksum_id_v16, p); err != nil {
		assert.NoError(t, err)
	}

	f, err := ioutil.ReadFile(p)
	require.NoError(t, err)
	snapshot.Assert(t, f)
}

func TestChecksumParse(t *testing.T) {
	data, err := ioutil.ReadFile("__snapshots__/testdownload.snap")
	require.NoError(t, err)

	checksum, err := getChecksum(data, Windows)
	assert.NoError(t, err)
	assert.Equal(t, checksum_v16_windows, checksum)
}

func TestUnarchiveAndChecksum(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Logf("Skipped expensive unarchive test. Set INTEGRATION_TESTS=true to enable.")
		t.SkipNow()
	}
	tmpdir, err := ioutil.TempDir("", "vat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	path := filepath.Join(tmpdir, "update.zip")
	if err := downloadAsset(nil, binary_windows_id_v16, path); err != nil {
		require.NoError(t, err)
	}

	checked, err := checksum(path, checksum_v16_windows)
	assert.NoError(t, err)
	assert.True(t, checked)

	if err := unarchive(path); err != nil {
		assert.NoError(t, err)
	}
	finfo, err := os.Stat(filepath.Join(tmpdir, "vat.exe"))
	assert.False(t, os.IsNotExist(err))
	assert.Equal(t, "vat.exe", finfo.Name())
}
