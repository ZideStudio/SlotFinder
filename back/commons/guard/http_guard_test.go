package guard

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func multipartRequest(t *testing.T, fieldSize int) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.bin")
	require.NoError(t, err)
	_, err = part.Write(make([]byte, fieldSize))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestMaxUploadSizeMiddleware_WithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = multipartRequest(t, 100)

	nextCalled := false
	MaxUploadSizeMiddleware(1 << 20)(c)
	if !c.IsAborted() {
		nextCalled = true
	}

	assert.True(t, nextCalled)
}

func TestMaxUploadSizeMiddleware_ExceedsLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = multipartRequest(t, 2000)

	MaxUploadSizeMiddleware(500)(c)

	assert.True(t, c.IsAborted())
}
