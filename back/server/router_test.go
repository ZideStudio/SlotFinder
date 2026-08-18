package server

import (
	"app/db"
	"app/testutils"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter_Healthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	req := httptest.NewRequest("GET", "/healthz", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "work", recorder.Body.String())
}

// Uses a dedicated connection, not the shared base: db.TestConnection()
// closes whatever pool it's handed, which would break tests running after.
func TestNewRouter_Readyz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	base := testutils.BaseDB(t)
	db.SetDBForTests(testutils.FreshDB(t))
	t.Cleanup(func() { db.SetDBForTests(base) })

	router := NewRouter()

	req := httptest.NewRequest("GET", "/readyz", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "ready", recorder.Body.String())
}

func TestNewRouter_Swagger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	req := httptest.NewRequest("GET", "/swagger/index.html", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
}

func TestNewRouter_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 404, recorder.Code)
}

func TestNewRouter_AccountCreateRequiresBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	req := httptest.NewRequest("POST", "/v1/account", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 400, recorder.Code)
}

func TestNewRouter_AccountMeRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	req := httptest.NewRequest("GET", "/v1/account/me", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 401, recorder.Code)
}
