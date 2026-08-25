package lib

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// RenderSVGToPNG rasterizes SVG markup onto an opaque white size×size canvas as PNG
func RenderSVGToPNG(svg []byte, size int) ([]byte, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, err
	}
	if icon.ViewBox.W <= 0 || icon.ViewBox.H <= 0 {
		return nil, fmt.Errorf("svg has no usable viewBox or width/height")
	}
	icon.SetTarget(0, 0, float64(size), float64(size))

	// oksvg doesn't scale stroke-width/dash with SetTarget, scale them by hand
	scale := float64(size) / icon.ViewBox.W
	for i := range icon.SVGPaths {
		p := &icon.SVGPaths[i]
		p.LineWidth *= scale
		p.DashOffset *= scale
		for j := range p.Dash {
			p.Dash[j] *= scale
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	raster := rasterx.NewDasher(size, size, scanner)
	icon.Draw(raster, 1.0)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
