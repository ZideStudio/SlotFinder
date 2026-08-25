package lib

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	return img
}

// isDark reports whether a pixel is closer to black than white, so a stroke
// of any color (as long as it contrasts with the white background) can be
// measured without hardcoding the exact SVG stroke color.
func isDark(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	return r+g+b < 3*0x8000
}

func TestRenderSVGToPNG_ValidSVG_ReturnsCorrectlySizedPNG(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><rect x="0" y="0" width="10" height="10" fill="#000000"/></svg>`

	data, err := RenderSVGToPNG([]byte(svg), 64)
	require.NoError(t, err)

	img := decodePNG(t, data)
	assert.Equal(t, 64, img.Bounds().Dx())
	assert.Equal(t, 64, img.Bounds().Dy())
}

// TestRenderSVGToPNG_ScalesStrokeWidth is a regression test for a bug where
// oksvg's SetTarget scales path geometry but not stroke-width, so a stroke
// rasterized at a size much larger than the SVG's viewBox came out far
// thinner than intended (reported as "avatars render thinner than the
// Dicebear playground"). A 1-unit-wide stroke in a 10x10 viewBox, rendered
// at 100x100 (10x scale), must come out roughly 10px wide, not 1px.
func TestRenderSVGToPNG_ScalesStrokeWidth(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><line x1="5" y1="1" x2="5" y2="9" stroke="#000000" stroke-width="1"/></svg>`

	data, err := RenderSVGToPNG([]byte(svg), 100)
	require.NoError(t, err)
	img := decodePNG(t, data)

	// Scan the horizontal midline for the run of dark pixels around the
	// vertical stroke at x=50. Unscaled, the stroke would be ~1px wide;
	// correctly scaled 10x, it should be close to 10px wide.
	run := 0
	for x := 0; x < 100; x++ {
		if isDark(img, x, 50) {
			run++
		}
	}
	assert.Greater(t, run, 6, "stroke should be scaled up with the rest of the geometry, not left at its raw SVG-unit width")
}

func TestRenderSVGToPNG_NoViewBoxOrSize_ReturnsError(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg"><circle cx="5" cy="5" r="5" fill="#000000"/></svg>`

	_, err := RenderSVGToPNG([]byte(svg), 64)
	assert.Error(t, err)
}

func TestRenderSVGToPNG_MalformedSVG_ReturnsError(t *testing.T) {
	_, err := RenderSVGToPNG([]byte(`<svg><circle cx="5"`), 64)
	assert.Error(t, err)
}
