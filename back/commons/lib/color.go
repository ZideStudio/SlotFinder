package lib

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
)

const (
	// dominantColorBucketShift keeps the 4 most significant bits of each channel so JPEG noise and smooth gradients don't split the majority color across many bins
	dominantColorBucketShift = 4
	// dominantColorMinAlpha is the alpha (8-bit) under which a pixel is considered transparent
	dominantColorMinAlpha = 128
	// dominantColorWhiteThreshold is the channel value from which all-channel pixels are considered white
	dominantColorWhiteThreshold = 245
)

var errNoDominantColor = errors.New("no dominant color: image is empty, transparent or white")

// DominantColor returns the most frequent color of a JPEG/PNG/WebP image as an uppercase "#RRGGBB" string
func DominantColor(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	type bucket struct {
		count            int
		sumR, sumG, sumB int
	}
	const bits = 8 - dominantColorBucketShift // bits kept per channel
	var buckets [1 << (3 * bits)]bucket

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if c.A < dominantColorMinAlpha {
				continue
			}
			if c.R >= dominantColorWhiteThreshold && c.G >= dominantColorWhiteThreshold && c.B >= dominantColorWhiteThreshold {
				continue
			}

			idx := int(c.R>>dominantColorBucketShift)<<(2*bits) |
				int(c.G>>dominantColorBucketShift)<<bits |
				int(c.B>>dominantColorBucketShift)
			buckets[idx].count++
			buckets[idx].sumR += int(c.R)
			buckets[idx].sumG += int(c.G)
			buckets[idx].sumB += int(c.B)
		}
	}

	best := -1
	for i := range buckets {
		if buckets[i].count > 0 && (best == -1 || buckets[i].count > buckets[best].count) {
			best = i
		}
	}
	if best == -1 {
		return "", errNoDominantColor
	}

	w := buckets[best]
	return fmt.Sprintf("#%02X%02X%02X", w.sumR/w.count, w.sumG/w.count, w.sumB/w.count), nil
}
