package guard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsCheck_MatchingOrigin_SetsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Origin", "https://slotfinder.test")

	g := &CorsGuard{}
	g.CorsCheck()(c)

	assert.Equal(t, "https://slotfinder.test", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCorsCheck_NonMatchingOrigin_NoHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Origin", "https://evil.example.com")

	g := &CorsGuard{}
	g.CorsCheck()(c)

	assert.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
}

func TestCorsCheck_OptionsRequest_Aborts204(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodOptions, "/", nil)

	g := &CorsGuard{}
	g.CorsCheck()(c)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}
