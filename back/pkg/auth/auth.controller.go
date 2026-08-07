package auth

import (
	"app/commons/helpers"
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	refreshTokenRepository *repository.RefreshTokenRepository
	cleanupCtx             context.Context
	cleanupCancel          context.CancelFunc
}

const refreshTokenCleanupInterval = 24 * time.Hour

// refreshCleanupClock is swapped for a fake clock in tests so the cleanup
// ticker can be advanced deterministically instead of waiting on the real
// 24h interval.
var refreshCleanupClock clock = realClock{}

func NewAuthController(ctl *AuthController) *AuthController {
	if ctl == nil {
		ctl = &AuthController{
			refreshTokenRepository: repository.NewRefreshTokenRepository(nil),
		}
	}

	ctl.cleanupCtx, ctl.cleanupCancel = context.WithCancel(context.Background())
	// Create the ticker synchronously here (rather than inside the goroutine)
	// so tests can swap refreshCleanupClock and advance it without racing the
	// background goroutine's first read of it.
	t := refreshCleanupClock.NewTicker(refreshTokenCleanupInterval)
	go ctl.cleanRefreshTokens(t)

	return ctl
}

// @Summary Status Check
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} helpers.ApiError
// @Failure 403 {object} helpers.ApiError
// @Router /api/v1/auth/status [get]
func (ctl *AuthController) Status(c *gin.Context) {
	helpers.HandleJSONResponse(c, nil, nil)
}

// @Summary Logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} helpers.ApiError
// @Router /api/v1/auth/logout [post]
// @security AccessTokenCookie
func (ctl *AuthController) Logout(c *gin.Context) {
	// Get the current refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		// Revoke only the current refresh token (logout from this device only)
		tokenHash := ctl.refreshTokenRepository.HashToken(refreshToken)
		var token model.RefreshToken
		if err := ctl.refreshTokenRepository.FindByTokenHash(tokenHash, &token); err == nil {
			_ = ctl.refreshTokenRepository.Revoke(token.Id)
		}
	}

	// Clear cookies
	lib.SetAccessTokenCookie(c, "", -1)
	lib.SetRefreshTokenCookie(c, "", -1)

	helpers.HandleJSONResponse(c, nil, nil)
}

// cleanRefreshTokens runs periodic refresh-token cleanup and stops when the controller is canceled.
func (ctl *AuthController) cleanRefreshTokens(t ticker) {
	defer t.Stop()

	for {
		select {
		case <-t.C():
			_ = ctl.refreshTokenRepository.DeleteExpired()
		case <-ctl.cleanupCtx.Done():
			return
		}
	}
}
