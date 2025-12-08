package signin

import (
	"app/commons/helpers"
	"app/commons/lib"

	"github.com/gin-gonic/gin"
)

type SigninController struct {
	signinService *SigninService
}

func NewSigninController(ctl *SigninController) *SigninController {
	if ctl != nil {
		return ctl
	}

	return &SigninController{
		signinService: NewSigninService(nil),
	}
}

// @Summary Sign in
// @Description Sign in with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body SigninDto true "Sign in parameters"
// @Success 200 {object} TokenResponseDto
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/auth/signin [post]
func (ctl *SigninController) Signin(c *gin.Context) {
	var data SigninDto
	err := helpers.SetHttpContextBody(c, &data)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	token, err := ctl.signinService.Signin(&data)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	lib.SetAccessTokenCookie(c, token.AccessToken, 0)
	lib.SetRefreshTokenCookie(c, token.RefreshToken, 0)

	helpers.HandleJSONResponse(c, nil, nil)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} helpers.ApiError
// @Router /api/v1/auth/refresh [post]
func (ctl *SigninController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	token, err := ctl.signinService.RefreshAccessToken(refreshToken)
	if err != nil {
		// Clear both cookies on error
		lib.SetAccessTokenCookie(c, "", -1)
		lib.SetRefreshTokenCookie(c, "", -1)
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	lib.SetAccessTokenCookie(c, token.AccessToken, 0)
	lib.SetRefreshTokenCookie(c, token.RefreshToken, 0)

	helpers.HandleJSONResponse(c, nil, nil)
}
