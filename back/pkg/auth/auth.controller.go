package auth

import (
	"app/commons/guard"
	"app/commons/helpers"
	"app/commons/lib"
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
	// Get user claims to revoke all refresh tokens
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err == nil && user != nil {
		// Revoke all refresh tokens for this user
		_ = ctl.refreshTokenRepository.RevokeAllForAccount(user.Id)
	}

	// Clear cookies
	lib.SetAccessTokenCookie(c, "", -1)
	lib.SetRefreshTokenCookie(c, "", -1)
	
	helpers.HandleJSONResponse(c, nil, nil)
}
