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

func NewAuthController(ctl *AuthController) *AuthController {
	if ctl == nil {
		ctl = &AuthController{
			refreshTokenRepository: &repository.RefreshTokenRepository{},
		}
	}

	ctl.cleanupCtx, ctl.cleanupCancel = context.WithCancel(context.Background())
	go ctl.cleanRefreshTokens()

	return ctl
}

// @Summary Status Check
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} helpers.ApiError
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

// cleanRefreshTokens avec gestion propre du cycle de vie
func (ctl *AuthController) cleanRefreshTokens() {
	ticker := time.NewTicker(24 * time.Hour)
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
