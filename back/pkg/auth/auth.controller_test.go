package auth

import (
	model "app/db/models"
	"app/db/repository"
	"app/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: Do not call NewAuthController(nil) with the real 24h ticker in most
// tests — it starts a background goroutine wired to real dependencies.
// Build the controller struct directly instead, same pattern used elsewhere
// in this codebase for unit tests.

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
	db := testutils.TestDB(t)
	ctl := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}
	c, recorder := newTestContext()

	ctl.Logout(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}

func TestLogout_InvalidCookie_StillClearsCookies(t *testing.T) {
	db := testutils.TestDB(t)
	ctl := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}
	c, recorder := newTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "unknown-token"})

	ctl.Logout(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestLogout_ValidCookie_RevokesToken(t *testing.T) {
	db := testutils.TestDB(t)
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
	ctl := NewAuthController(nil)
	defer ctl.cleanupCancel()

	assert.NotNil(t, ctl.refreshTokenRepository)
	assert.NotNil(t, ctl.cleanupCtx)
	assert.NotNil(t, ctl.cleanupCancel)
}

func TestNewAuthController_ReusesProvidedInstance(t *testing.T) {
	db := testutils.TestDB(t)
	provided := &AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)}

	ctl := NewAuthController(provided)
	defer ctl.cleanupCancel()

	assert.Same(t, provided, ctl)
}

func TestCleanRefreshTokens_DeletesExpiredOnTick(t *testing.T) {
	// Use the real 24h production interval, but drive it with a fake clock
	// advanced past that interval so the test doesn't wait on wall-clock time.
	original := refreshCleanupClock
	fake := newFakeClock()
	refreshCleanupClock = fake
	defer func() { refreshCleanupClock = original }()

	db := testutils.TestDB(t)
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

	fake.Advance(refreshTokenCleanupInterval)

	// A single wait-then-check instead of require.Eventually: Eventually
	// polls by spawning a new goroutine per tick, and every one of them would
	// share this test's single dedicated Postgres connection (testutils.TestDB)
	// with the cleanup goroutine's own DeleteExpired() call — concurrent
	// access jackc/pgx connections don't support and that intermittently
	// fails with "driver: bad connection".
	testutils.AwaitAsyncDBWork(t)

	var count int64
	require.NoError(t, db.Model(&model.RefreshToken{}).Where("id = ?", expired.Id).Count(&count).Error)
	assert.Equal(t, int64(0), count, "expected the cleanup ticker to delete the expired token")
}

func TestCleanRefreshTokens_StopsOnCancel(t *testing.T) {
	// The real 24h interval never fires within this test's lifetime, so the
	// only thing to verify is that cancellation stops the goroutine cleanly.
	db := testutils.TestDB(t)
	ctl := NewAuthController(&AuthController{refreshTokenRepository: repository.NewRefreshTokenRepository(db)})

	ctl.cleanupCancel()

	// Give the goroutine a moment to observe cancellation and return; there's
	// nothing further to assert beyond "this doesn't hang/panic".
	time.Sleep(50 * time.Millisecond)
}

func TestFakeTicker_StoppedTicker_DoesNotFire(t *testing.T) {
	fake := newFakeClock()
	tk := fake.NewTicker(time.Minute)
	tk.Stop()

	fake.Advance(time.Minute)

	select {
	case <-tk.C():
		t.Fatal("expected no tick after Stop")
	default:
	}
}
