package provider

import (
	"app/commons/encryption"
	"app/commons/guard"
	"app/commons/helpers"
	"encoding/base64"
	"encoding/json"
	"errors"

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
// @Success 200 {string} string "OAuth URL"
// @Failure 400 {object} helpers.ApiError
// @Router /v1/auth/{provider}/url [get]
func (ctl *ProviderController) ProviderUrl(c *gin.Context) {
	provider := c.Param("provider")
	redirectUrl := c.Query("redirectUrl")

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.ResponseJSON(c, nil, err)
		return
	}

	url, err := ctl.signinService.GetProviderUrl(provider, redirectUrl, user)

	helpers.ResponseJSON(c, url, err)
}

func (ctl *ProviderController) ProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	// get state
	stateStrEncrypted := c.Query("state")
	if stateStrEncrypted == "" {
		helpers.ResponseJSON(c, nil, errors.New("State parameter is missing"))
		return
	}
	// decrypt state
	decodedState, err := encryption.Decrypt(stateStrEncrypted)
	if err != nil {
		helpers.ResponseJSON(c, nil, errors.New("State parameter has invalid encryption"))
		return
	}

	var state map[string]string
	json.Unmarshal([]byte(decodedState), &state)
	redirectUrl := state["redirectUrl"]
	if state == nil || redirectUrl == "" {
		helpers.ResponseJSON(c, nil, errors.New("State parameter is invalid"))
		return
	}

	userId := state["userId"]

	jwt, err := ctl.signinService.ProviderCallback(provider, code, userId)

	if err != nil {
		c.Redirect(302, redirectUrl+"?error="+base64.StdEncoding.EncodeToString([]byte(err.Error())))
		return
	}

	c.Redirect(302, redirectUrl+"?access_token="+base64.StdEncoding.EncodeToString([]byte(jwt.AccessToken)))
}
