package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCookieTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, recorder
}

func TestSetAccessTokenCookie_DefaultExpiration(t *testing.T) {
	t.Parallel()
	c, recorder := newCookieTestContext()
	SetAccessTokenCookie(c, "token-value", 0)

	cookies := recorder.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "access_token", cookies[0].Name)
	assert.Equal(t, "token-value", cookies[0].Value)
	assert.Equal(t, "/api", cookies[0].Path)
	assert.True(t, cookies[0].MaxAge > 0)
}

func TestSetAccessTokenCookie_ExplicitExpiration(t *testing.T) {
	t.Parallel()
	c, recorder := newCookieTestContext()
	SetAccessTokenCookie(c, "token-value", -1)

	cookies := recorder.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, -1, cookies[0].MaxAge)
}

func TestSetRefreshTokenCookie_DefaultExpiration(t *testing.T) {
	t.Parallel()
	c, recorder := newCookieTestContext()
	SetRefreshTokenCookie(c, "refresh-value", 0)

	cookies := recorder.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "refresh_token", cookies[0].Name)
	assert.Equal(t, "/api/v1/auth/refresh", cookies[0].Path)
	assert.True(t, cookies[0].MaxAge > 0)
}

func TestSetRefreshTokenCookie_ExplicitExpiration(t *testing.T) {
	t.Parallel()
	c, recorder := newCookieTestContext()
	SetRefreshTokenCookie(c, "", -1)

	cookies := recorder.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, -1, cookies[0].MaxAge)
}
