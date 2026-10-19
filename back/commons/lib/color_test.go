package lib

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// encodeStripedPNG builds a w×h PNG whose first `firstCols` columns are `first` and the rest `second`.
func encodeStripedPNG(t *testing.T, w, h, firstCols int, first, second color.NRGBA) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if x < firstCols {
				img.SetNRGBA(x, y, first)
			} else {
				img.SetNRGBA(x, y, second)
			}
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func encodeSolidPNG(t *testing.T, c color.NRGBA) []byte {
	t.Helper()
	return encodeStripedPNG(t, 10, 10, 10, c, c)
}

var (
	testRed         = color.NRGBA{R: 255, A: 255}
	testBlue        = color.NRGBA{B: 255, A: 255}
	testGreen       = color.NRGBA{G: 200, A: 255}
	testWhite       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	testTransparent = color.NRGBA{}
)

func TestDominantColor_SolidColor(t *testing.T) {
	got, err := DominantColor(encodeSolidPNG(t, testRed))

	require.NoError(t, err)
	assert.Equal(t, "#FF0000", got)
}

func TestDominantColor_MajorityWins(t *testing.T) {
	// 7 red columns vs 3 blue columns.
	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 7, testRed, testBlue))

	require.NoError(t, err)
	assert.Equal(t, "#FF0000", got)
}

func TestDominantColor_IgnoresWhitePixels(t *testing.T) {
	// 8 white columns vs 2 green columns.
	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 8, testWhite, testGreen))

	require.NoError(t, err)
	assert.Equal(t, "#00C800", got)
}

func TestDominantColor_IgnoresNearWhitePixels(t *testing.T) {
	nearWhite := color.NRGBA{R: 247, G: 248, B: 250, A: 255}

	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 8, nearWhite, testGreen))

	require.NoError(t, err)
	assert.Equal(t, "#00C800", got)
}

func TestDominantColor_IgnoresTransparentPixels(t *testing.T) {
	// 8 transparent columns vs 2 blue columns.
	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 8, testTransparent, testBlue))

	require.NoError(t, err)
	assert.Equal(t, "#0000FF", got)
}

func TestDominantColor_AveragesPixelsOfWinningBucket(t *testing.T) {
	// 0x80 and 0x8F share the same bucket (high nibble 8): the result is their average, not the bucket center.
	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 5, color.NRGBA{R: 0x80, A: 255}, color.NRGBA{R: 0x8F, A: 255}))

	require.NoError(t, err)
	assert.Equal(t, "#870000", got)
}

func TestDominantColor_TieResolvedByLowestBucket(t *testing.T) {
	// 5 red columns vs 5 blue columns: blue has the lowest bucket index (0x00F vs 0xF00).
	got, err := DominantColor(encodeStripedPNG(t, 10, 10, 5, testRed, testBlue))

	require.NoError(t, err)
	assert.Equal(t, "#0000FF", got)
}

func TestDominantColor_OutputIsValidHex(t *testing.T) {
	got, err := DominantColor(encodeSolidPNG(t, color.NRGBA{R: 12, G: 34, B: 56, A: 255}))

	require.NoError(t, err)
	assert.True(t, IsHexa(got))
	assert.Len(t, got, 7)
}

func TestDominantColor_JPEGInput(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			img.SetNRGBA(x, y, testBlue)
		}
	}
	var buf bytes.Buffer
	require.NoError(t, jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}))

	got, err := DominantColor(buf.Bytes())

	require.NoError(t, err)
	// JPEG compression is lossy: accept a small deviation on each channel.
	assert.Len(t, got, 7)
	r, g, b := parseHexColor(t, got)
	assert.InDelta(t, 0, r, 20)
	assert.InDelta(t, 0, g, 20)
	assert.InDelta(t, 255, b, 20)
}

func TestDominantColor_AllWhite(t *testing.T) {
	got, err := DominantColor(encodeSolidPNG(t, testWhite))

	assert.ErrorIs(t, err, errNoDominantColor)
	assert.Empty(t, got)
}

func TestDominantColor_AllTransparent(t *testing.T) {
	got, err := DominantColor(encodeSolidPNG(t, testTransparent))

	assert.ErrorIs(t, err, errNoDominantColor)
	assert.Empty(t, got)
}

func TestDominantColor_InvalidData(t *testing.T) {
	got, err := DominantColor([]byte("not an image"))

	assert.Error(t, err)
	assert.Empty(t, got)
}

func TestDominantColor_EmptyData(t *testing.T) {
	got, err := DominantColor(nil)

	assert.Error(t, err)
	assert.Empty(t, got)
}

func parseHexColor(t *testing.T, hex string) (r, g, b int) {
	t.Helper()
	_, err := fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	require.NoError(t, err)
	return r, g, b
}
