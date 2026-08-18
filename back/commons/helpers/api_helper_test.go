package helpers

import (
	"app/commons/constants"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHelperTestContext(body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, recorder
}

type sampleBody struct {
	Name string `json:"name" binding:"required"`
}

func TestSetHttpContextBody_Success(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext([]byte(`{"name":"hello"}`))

	var body sampleBody
	err := SetHttpContextBody(c, &body)
	assert.NoError(t, err)
	assert.Equal(t, "hello", body.Name)
	assert.Equal(t, http.StatusOK, recorder.Code) // nothing written on success
}

func TestSetHttpContextBody_InvalidJSON(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext([]byte(`not-json`))

	var body sampleBody
	err := SetHttpContextBody(c, &body)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSetHttpContextBody_MissingRequiredField(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext([]byte(`{}`))

	var body sampleBody
	err := SetHttpContextBody(c, &body)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandleJSONResponse_Success(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext(nil)

	HandleJSONResponse(c, map[string]string{"key": "value"}, nil)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var parsed map[string]string
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &parsed))
	assert.Equal(t, "value", parsed["key"])
}

func TestHandleJSONResponse_CustomError_WithStatus(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext(nil)

	HandleJSONResponse(c, nil, constants.ERR_NOT_AUTHENTICATED.Err)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	var parsed ApiError
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &parsed))
	assert.Equal(t, "NOT_AUTHENTICATED", parsed.Code)
}

func TestHandleJSONResponse_CustomError_DefaultStatus(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext(nil)

	HandleJSONResponse(c, nil, constants.ERR_INVALID_EMAIL_FORMAT.Err)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandleJSONResponse_UnknownError_ServerError(t *testing.T) {
	t.Parallel()
	c, recorder := newHelperTestContext(nil)

	HandleJSONResponse(c, nil, errors.New("something unexpected"))

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	var parsed ApiError
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &parsed))
	assert.Equal(t, "SERVER_ERROR", parsed.Code)
}
