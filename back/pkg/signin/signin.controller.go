package signin

import (
	"app/commons/helpers"

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
	helpers.ShouldBindJSON(c, &data)

	token, err := ctl.signinService.Signin(&data)

	helpers.ResponseJSON(c, token, err)
}
