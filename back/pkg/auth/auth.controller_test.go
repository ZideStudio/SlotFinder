package auth

import (
	model "app/db/models"
	"app/db/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NOTE: Do not call NewAuthController(nil) with the real 24h ticker in most
// tests — it starts a background goroutine wired to real dependencies.
// Build the controller struct directly instead, same pattern used elsewhere
// in this codebase for unit tests.

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.Account{}, &model.RefreshToken{}))

	// A ":memory:" sqlite DSN is a fresh, isolated DB per *connection*. Some
	// of these tests query it concurrently from a background goroutine (the
	// cleanup ticker), so pin the pool to a single connection or those
	// queries would silently hit an empty, unmigrated database.
	sqlDB, err := database.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	return database
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	return c, recorder
}

func TestStatus(t *testing.T) {
	ctl := &AuthController{}
	c, recorder := newTestContext()

	ctl.Status(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestLogout_NoCookie(t *testing.T) {
	db := testDB(t)
	ctl := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}
	c, recorder := newTestContext()

	ctl.Logout(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}

func TestLogout_InvalidCookie_StillClearsCookies(t *testing.T) {
	db := testDB(t)
	ctl := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}
	c, recorder := newTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "unknown-token"})

	ctl.Logout(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestLogout_ValidCookie_RevokesToken(t *testing.T) {
	db := testDB(t)
	repo := repository.NewRefreshTokenRepository(db)
	ctl := &AuthController{refreshTokenRepository: repo}

	account := model.Account{Id: uuid.New()}
	require.NoError(t, db.Create(&account).Error)
	token, err := repo.Create(account.Id, time.Now().Add(time.Hour))
	require.NoError(t, err)

	c, recorder := newTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: token})

	ctl.Logout(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var stored model.RefreshToken
	require.NoError(t, db.Where("token_hash = ?", repo.HashToken(token)).First(&stored).Error)
	assert.True(t, stored.IsRevoked)
}

func TestNewAuthController_Nil_BuildsDefaultAndStartsCleanup(t *testing.T) {
	original := refreshTokenCleanupInterval
	refreshTokenCleanupInterval = time.Hour // avoid firing during this fast test
	defer func() { refreshTokenCleanupInterval = original }()

	ctl := NewAuthController(nil)
	defer ctl.cleanupCancel()

	assert.NotNil(t, ctl.refreshTokenRepository)
	assert.NotNil(t, ctl.cleanupCtx)
	assert.NotNil(t, ctl.cleanupCancel)
}

func TestNewAuthController_ReusesProvidedInstance(t *testing.T) {
	original := refreshTokenCleanupInterval
	refreshTokenCleanupInterval = time.Hour
	defer func() { refreshTokenCleanupInterval = original }()

	db := testDB(t)
	provided := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}

	ctl := NewAuthController(provided)
	defer ctl.cleanupCancel()

	assert.Same(t, provided, ctl)
}

func TestCleanRefreshTokens_DeletesExpiredOnTick(t *testing.T) {
	original := refreshTokenCleanupInterval
	refreshTokenCleanupInterval = 10 * time.Millisecond
	defer func() { refreshTokenCleanupInterval = original }()

	db := testDB(t)
	repo := repository.NewRefreshTokenRepository(db)
	account := model.Account{Id: uuid.New()}
	require.NoError(t, db.Create(&account).Error)

	expired := model.RefreshToken{
		Id:        uuid.New(),
		AccountId: account.Id,
		TokenHash: "expired",
		ExpiresAt: time.Now().AddDate(0, 0, -10),
		IsRevoked: true,
	}
	require.NoError(t, db.Create(&expired).Error)

	ctl := NewAuthController(&AuthController{refreshTokenRepository: repo})
	defer ctl.cleanupCancel()

	require.Eventually(t, func() bool {
		var count int64
		db.Model(&model.RefreshToken{}).Where("id = ?", expired.Id).Count(&count)
		return count == 0
	}, time.Second, 10*time.Millisecond, "expected the cleanup ticker to delete the expired token")
}

func TestCleanRefreshTokens_StopsOnCancel(t *testing.T) {
	original := refreshTokenCleanupInterval
	refreshTokenCleanupInterval = 5 * time.Millisecond
	defer func() { refreshTokenCleanupInterval = original }()

	db := testDB(t)
	ctl := NewAuthController(&AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)})

	ctl.cleanupCancel()

	// Give the goroutine a moment to observe cancellation and return; there's
	// nothing further to assert beyond "this doesn't hang/panic".
	time.Sleep(50 * time.Millisecond)
}
