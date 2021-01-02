package ocr

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/BTBurke/vatinator/img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: make this more robust beyond just smoke test
func TestProcess(t *testing.T) {
	snap, err := snapshot.New(snapshot.SnapExtension(".json"))
	require.NoError(t, err)

	fname := "../test_receipts/PXL_20201002_163306793.jpg"
	f, err := os.Open(fname)
	require.NoError(t, err)
	image, err := img.NewImageFromReader(f)
	require.NoError(t, err)

	res, err := ProcessImage(image, "../vatinator-f91ccb107c2c.json")
	require.NoError(t, err)

	resB, err := json.Marshal(res)
	assert.NoError(t, err)
	snap.Assert(t, resB)
}

func TestJoin(t *testing.T) {
	in := []string{"this", "is", "a", "test"}
	out := append(in, "this is", "is a", "a test")
	assert.Equal(t, out, joinFollowing(in))
}
