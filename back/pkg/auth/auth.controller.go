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

// refreshTokenCleanupInterval is overridable in tests to avoid waiting 24h
// for the cleanup ticker to fire.
var refreshTokenCleanupInterval = 24 * time.Hour

func NewAuthController(ctl *AuthController) *AuthController {
	if ctl == nil {
		ctl = &AuthController{
			refreshTokenRepository: repository.NewRefreshTokenRepository(nil),
		}
	}

	ctl.cleanupCtx, ctl.cleanupCancel = context.WithCancel(context.Background())
	// Read the interval synchronously here (rather than inside the goroutine)
	// so tests can override refreshTokenCleanupInterval without racing with
	// the background goroutine's first read of it.
	interval := refreshTokenCleanupInterval
	go ctl.cleanRefreshTokens(interval)

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
func (ctl *AuthController) cleanRefreshTokens(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = ctl.refreshTokenRepository.DeleteExpired()
		case <-ctl.cleanupCtx.Done():
			return
		}
	}
}
