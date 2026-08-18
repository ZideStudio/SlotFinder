package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPagination_ParseQuery_NilContext(t *testing.T) {
	t.Parallel()
	var p Pagination[string]
	err := p.ParseQuery(nil)
	assert.Error(t, err)
}

func TestPagination_ParseQuery_Defaults(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	var p Pagination[string]
	err := p.ParseQuery(c)
	assert.NoError(t, err)
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Limit)
	assert.Equal(t, 0, p.Offset)
}

func TestPagination_ParseQuery_ComputesOffset(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=3&limit=10", nil)

	var p Pagination[string]
	err := p.ParseQuery(c)
	assert.NoError(t, err)
	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 10, p.Limit)
	assert.Equal(t, 20, p.Offset)
}

func TestPagination_ParseQuery_InvalidParams(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=0", nil)

	var p Pagination[string]
	err := p.ParseQuery(c)
	assert.Error(t, err)
}
