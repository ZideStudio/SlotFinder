package lib

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func decodeJPEGSize(t *testing.T, data []byte) (int, int) {
	t.Helper()
	img, err := jpeg.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	b := img.Bounds()
	return b.Dx(), b.Dy()
}

func TestProcessAvatar_SmallImage_KeepsDimensions(t *testing.T) {
	input := buildPNG(t, 100, 50)
	out, err := ProcessAvatar(input)
	require.NoError(t, err)

	w, h := decodeJPEGSize(t, out)
	assert.Equal(t, 100, w)
	assert.Equal(t, 50, h)
}

func TestProcessAvatar_LargeImage_ResizesToMaxDimension(t *testing.T) {
	input := buildPNG(t, 1000, 500) // 2:1 ratio, width is the larger side
	out, err := ProcessAvatar(input)
	require.NoError(t, err)

	w, h := decodeJPEGSize(t, out)
	assert.Equal(t, 256, w)
	assert.Equal(t, 128, h) // ratio preserved
}

func TestProcessAvatar_LargeImage_TallAspectRatio(t *testing.T) {
	input := buildPNG(t, 500, 1000) // height is the larger side
	out, err := ProcessAvatar(input)
	require.NoError(t, err)

	w, h := decodeJPEGSize(t, out)
	assert.Equal(t, 128, w)
	assert.Equal(t, 256, h)
}

func TestProcessAvatar_OutputIsValidJPEG(t *testing.T) {
	input := buildPNG(t, 20, 20)
	out, err := ProcessAvatar(input)
	require.NoError(t, err)
	assert.NotZero(t, len(out))

	_, err = jpeg.Decode(bytes.NewReader(out))
	assert.NoError(t, err, "output should be a valid JPEG")
}

func TestProcessAvatar_ExtremeWideAspectRatio_ClampsHeightTo1(t *testing.T) {
	input := buildPNG(t, 1000, 1)
	out, err := ProcessAvatar(input)
	require.NoError(t, err)

	w, h := decodeJPEGSize(t, out)
	assert.Equal(t, 256, w)
	assert.Equal(t, 1, h)
}

func TestProcessAvatar_ExtremeTallAspectRatio_ClampsWidthTo1(t *testing.T) {
	input := buildPNG(t, 1, 1000)
	out, err := ProcessAvatar(input)
	require.NoError(t, err)

	w, h := decodeJPEGSize(t, out)
	assert.Equal(t, 1, w)
	assert.Equal(t, 256, h)
}

func TestProcessAvatar_InvalidBytes(t *testing.T) {
	_, err := ProcessAvatar([]byte("not-an-image"))
	assert.Error(t, err)
}

func TestProcessAvatarFromURL_Success(t *testing.T) {
	body := buildPNG(t, 20, 20)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer server.Close()

	out, err := ProcessAvatarFromURL(server.URL)
	require.NoError(t, err)
	_, err = jpeg.Decode(bytes.NewReader(out))
	assert.NoError(t, err)
}

func TestProcessAvatarFromURL_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := ProcessAvatarFromURL(server.URL)
	assert.Error(t, err)
}

func TestProcessAvatarFromURL_ConnectionRefused(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	_, err := ProcessAvatarFromURL(url)
	assert.Error(t, err)
}

func TestProcessAvatarFromURL_ExceedsMaxDownloadBytes(t *testing.T) {
	chunk := bytes.Repeat([]byte{0}, 1<<20) // 1 MB chunk
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Stream slightly more than maxDownloadBytes (10 MB) of zero bytes.
		// Content doesn't need to be a valid image: the size check happens
		// before any decoding is attempted.
		for i := 0; i < 11; i++ {
			_, _ = w.Write(chunk)
		}
	}))
	defer server.Close()

	_, err := ProcessAvatarFromURL(server.URL)
	assert.Error(t, err)
}

func TestProcessAvatarFromURL_BodyReadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Declare more bytes than sent, then close the connection early to force a read error.
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("short"))
		hj, ok := w.(http.Hijacker)
		require.True(t, ok)
		conn, _, err := hj.Hijack()
		require.NoError(t, err)
		_ = conn.Close()
	}))
	defer server.Close()

	_, err := ProcessAvatarFromURL(server.URL)
	assert.Error(t, err)
}

func TestProcessAvatarFromURL_InvalidBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-an-image"))
	}))
	defer server.Close()

	_, err := ProcessAvatarFromURL(server.URL)
	assert.Error(t, err)
}
