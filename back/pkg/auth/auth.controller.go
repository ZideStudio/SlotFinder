package auth

import (
	"app/commons/helpers"
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	refreshTokenRepository *repository.RefreshTokenRepository
}

func NewAuthController() *AuthController {
	return &AuthController{
		refreshTokenRepository: &repository.RefreshTokenRepository{},
	}
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
