package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// NOTE: Ready()'s "db is up" branch (db.TestConnection() == true) is not
// exercised here: it depends on the process-global, unexported db.conn, which
// is only initialized by db.Init() against a real Postgres instance. That
// branch is covered by integration/E2E testing instead. TestConnection()
// itself has full unit coverage in app/db (db_test.go).

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, recorder
}

func TestHealthController_Status(t *testing.T) {
	ctl := HealthController{}
	c, recorder := newTestContext()

	ctl.Status(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "work", recorder.Body.String())
}

func TestHealthController_Ready_DbNotReady(t *testing.T) {
	ctl := HealthController{}
	c, recorder := newTestContext()

	ctl.Ready(c)

	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	assert.Equal(t, "db not ready", recorder.Body.String())
}
