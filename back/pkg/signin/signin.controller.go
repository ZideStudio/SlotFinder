package signin

import (
	"app/commons/helpers"
	"app/config"
	"time"

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
// @Router /v1/auth/signin [post]
func (ctl *SigninController) Signin(c *gin.Context) {
	var data SigninDto
	helpers.SetHttpContextBody(c, &data)

	token, err := ctl.signinService.Signin(&data)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	c.SetCookie(
		"access_token",            // name
		token.AccessToken,         // value
		int(168*time.Hour),        // max age in seconds
		"/",                       // path
		config.GetConfig().Domain, // domain
		true,                      // secure
		true,                      // httpOnly
	)

	helpers.HandleJSONResponse(c, nil, nil)
}
