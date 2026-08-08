package health

import (
	"app/db"
	model "app/db/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestHealthController_Ready_DbReady(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.Account{}))
	db.SetDBForTests(database)
	t.Cleanup(func() { db.SetDBForTests(nil) })

	ctl := HealthController{}
	c, recorder := newTestContext()

	ctl.Ready(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "ready", recorder.Body.String())
}
