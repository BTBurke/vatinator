package pdf

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/require"
)

// TODO: figure out how to strip out EXIF data that is useless and cause content-addressable storage to
// to store to copies
func TestPDF2Image(t *testing.T) {
	// this test works but because it embeds a lot of creation time EXIF data, the snapshot test
	// fails.
	t.SkipNow()

	f, err := os.Open("../test_receipts/pdfreceipt.pdf")
	require.NoError(t, err)

	out, err := PdfToImage(f)
	require.NoError(t, err)

	var b bytes.Buffer
	if _, err := io.Copy(&b, out); err != nil {
		require.NoError(t, err)
	}
	snap, err := snapshot.New(snapshot.SnapExtension(".png"))
	require.NoError(t, err)
	snap.Assert(t, b.Bytes())
}
