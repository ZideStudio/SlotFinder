package provider

import (
	"app/commons/constants"
	"app/commons/encryption"
	"app/commons/guard"
	"app/commons/helpers"
	"app/commons/lib"
	"app/config"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ProviderController struct {
	signinService *ProviderService
}

func NewProviderController(ctl *ProviderController) *ProviderController {
	if ctl != nil {
		return ctl
	}

	return &ProviderController{
		signinService: NewProviderService(nil),
	}
}

// @Summary Get redirect URL for OAuth provider
// @Tags Authentication
// @Param provider path string true "OAuth provider" Enums(google, github, discord)
// @Param redirectUrl query string true "URL to redirect after OAuth authentication"
// @Success 200 {string} string "OAuth URL"
// @Failure 400 {object} helpers.ApiError
// @Router /v1/auth/{provider}/url [get]
func (ctl *ProviderController) ProviderUrl(c *gin.Context) {
	provider := c.Param("provider")
	returnUrl := c.Query("returnUrl")

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	url, err := ctl.signinService.GetProviderUrl(provider, returnUrl, user)
	if err != nil {
		helpers.HandleJSONResponse(c, url, err)
		return
	}

	c.Redirect(302, url)
}

func (ctl *ProviderController) ProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	// get state
	stateStrEncrypted := c.Query("state")
	if stateStrEncrypted == "" {
		helpers.HandleJSONResponse(c, nil, errors.New("State parameter is missing"))
		return
	}
	// decrypt state
	decodedState, err := encryption.Decrypt(stateStrEncrypted)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, errors.New("State parameter has invalid encryption"))
		return
	}

	var state map[string]string
	if err := json.Unmarshal([]byte(decodedState), &state); err != nil {
		helpers.HandleJSONResponse(c, nil, errors.New("State parameter contains invalid JSON"))
		return
	} else if len(state) == 0 {
		helpers.HandleJSONResponse(c, nil, errors.New("State parameter is invalid"))
		return
	}

	redirectUrl := config.GetConfig().Origins[0] + "/oauth/callback"
	userId := state["userId"]
	returnUrl := state["returnUrl"]

	jwt, err := ctl.signinService.ProviderCallback(provider, code, userId)
	if err != nil {
		c.Redirect(302, fmt.Sprintf("%s?error=%s&returnUrl=%s", redirectUrl, constants.ERR_PROVIDER_CONNECTION_FAILED.Err.Error(), returnUrl))
		return
	}

	lib.SetAccessTokenCookie(c, jwt.AccessToken, 0)

	c.Redirect(302, fmt.Sprintf("%s?returnUrl=%s", redirectUrl, returnUrl))
}
