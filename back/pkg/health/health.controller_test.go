package health

import (
	"app/db"
	"app/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

// Uses a dedicated connection, not the shared base: db.TestConnection()
// closes whatever pool it's handed, which would break tests running after.
func TestHealthController_Ready_DbReady(t *testing.T) {
	database := testutils.FreshDB(t)
	db.SetDBForTests(database)
	t.Cleanup(func() { db.SetDBForTests(nil) })

	ctl := HealthController{}
	c, recorder := newTestContext()

	ctl.Ready(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "ready", recorder.Body.String())
}
