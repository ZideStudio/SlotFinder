package lib

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const avatarMaxDimension = 256

// ProcessAvatar decodes a JPEG/PNG/WebP image, resizes it to 256×256 max, and re-encodes it as JPEG 80%.
func ProcessAvatar(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	img = resizeIfNeeded(img, avatarMaxDimension)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ProcessAvatarFromURL fetches an image from a URL and runs it through ProcessAvatar.
func ProcessAvatarFromURL(imgURL string) ([]byte, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ProcessAvatar(data)
}

func resizeIfNeeded(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= maxDim && h <= maxDim {
		return img
	}

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = h * maxDim / w
	} else {
		newH = maxDim
		newW = w * maxDim / h
	}
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}
