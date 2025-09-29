package auth

import (
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// @Summary Status Check
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} helpers.ApiError
// @Router /v1/auth/status [get]
func (ctl *AuthController) Status(c *gin.Context) {
	helpers.HandleJSONResponse(c, nil, nil)
}
