package ocr

import (
	"context"
	"errors"
	"os"
	"testing"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/BTBurke/snapshot"
	"github.com/BTBurke/vatinator/img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func TestDetectOrientation(t *testing.T) {
	tt := []struct {
		name   string
		file   string
		expect Orientation
	}{
		{name: "rotate 0", file: "../test_receipts/PXL_20201002_163142937.jpg", expect: Orientation0},
		{name: "rotate 270", file: "../test_receipts/exif-6.jpg", expect: Orientation270},
		{name: "rotate 90", file: "../test_receipts/exif-8.jpg", expect: Orientation90},
	}
	snap, err := snapshot.New(snapshot.Diffable(false), snapshot.SnapExtension(".png"))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			imgReader, err := os.Open(tc.file)
			require.NoError(t, err)
			defer imgReader.Close()

			original, err := img.NewImageFromReader(imgReader)
			require.NoError(t, err)

			r, err := original.NewReader()
			require.NoError(t, err)
			i, err := vision.NewImageFromReader(r)
			require.NoError(t, err)
			ctx := context.Background()

			if _, err := os.Stat("../vatinator-f91ccb107c2c.json"); errors.Is(err, os.ErrNotExist) {
				t.Skipf("Skipping external API call: no GCS credentials")
			}

			c, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("../vatinator-f91ccb107c2c.json"))
			require.NoError(t, err)

			res, err := c.DetectTexts(ctx, i, &pb.ImageContext{LanguageHints: []string{"ET"}}, 1000)
			require.NoError(t, err)
			require.NotNil(t, res)

			orient := DetectOrientation(res)
			assert.Equal(t, tc.expect, DetectOrientation(res))
			if tc.expect != Orientation0 {
				rotated, err := AutoRotateImage(original, orient)
				require.NoError(t, err)
				png, err := rotated.AsPNG()
				require.NoError(t, err)
				snap.Assert(t, png)
			}

		})
	}
}
