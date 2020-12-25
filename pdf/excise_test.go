package pdf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/BTBurke/vatinator/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFillExciseForm(t *testing.T) {
	testExcise := types.Excise{
		Type:    "Gas 95",
		Arve:    "11111",
		Content: "test",
		Tax:     1000,
		Amount:  "25L",
		Date:    "24/12/2020",
	}

	testMD := types.ExciseMetadata{
		Embassy: "US Embassy",
		Name:    "Bryan Burke",
		Bank:    "My Bank, Liiviala 9, Tallinn 001110, EE220000000000",
		Date:    "December 2020",
	}
	testReceipts := []types.Excise{testExcise, testExcise, testExcise, testExcise, testExcise, testExcise}

	tmpdir, err := ioutil.TempDir("", "excise-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	outPath := filepath.Join(tmpdir, "out.pdf")
	if err := FillExcise(outPath, testReceipts, testMD, &FillExciseOptions{
		DisableRemote: true,
	}); err != nil {
		assert.NoError(t, err)
	}
	out, err := ioutil.ReadFile(outPath)
	require.NoError(t, err)
	snap, err := snapshot.New(snapshot.SnapExtension(".pdf"), snapshot.Diffable(false))
	require.NoError(t, err)
	snap.Assert(t, out)
}

func TestRemoteFill(t *testing.T) {
	testExcise := types.Excise{
		Type:    "Gas 95",
		Arve:    "11111",
		Content: "test",
		Tax:     1000,
		Amount:  "25L",
		Date:    "24/12/2020",
	}

	testMD := types.ExciseMetadata{
		Embassy: "US Embassy",
		Name:    "Bryan Burke",
		Bank:    "My Bank, Liiviala 9, Tallinn 001110, EE220000000000",
		Date:    "December 2020",
	}
	testReceipts := []types.Excise{testExcise, testExcise, testExcise, testExcise, testExcise, testExcise}

	tmpdir, err := ioutil.TempDir("", "excise-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	outPath := filepath.Join(tmpdir, "out.pdf")
	key := apiKey("../vatinator-f91ccb107c2c.json")
	require.Greater(t, len(key), 0)

	if err := FillExcise(outPath, testReceipts, testMD, &FillExciseOptions{
		ForceRemote: true,
		RemoteURL:   DefaultURL,
		APIKey:      apiKey("../vatinator-f91ccb107c2c.json"),
	}); err != nil {
		require.NoError(t, err)
	}
	out, err := ioutil.ReadFile(outPath)
	require.NoError(t, err)
	snap, err := snapshot.New(snapshot.SnapExtension(".pdf"), snapshot.Diffable(false))
	require.NoError(t, err)
	snap.Assert(t, out)
}
