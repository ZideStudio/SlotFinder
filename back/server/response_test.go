package server

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPRes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	HTTPRes(c, 200, "ok", map[string]string{"foo": "bar"})

	assert.Equal(t, 200, recorder.Code)
	assert.JSONEq(t, `{"code":200,"msg":"ok","data":{"foo":"bar"}}`, recorder.Body.String())
}

func TestHTTPRes_NilData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	HTTPRes(c, 404, "not found", nil)

	assert.Equal(t, 404, recorder.Code)
	assert.JSONEq(t, `{"code":404,"msg":"not found","data":null}`, recorder.Body.String())
}
