package provider

import (
	"app/commons/encryption"
	"app/commons/guard"
	"app/commons/helpers"
	"app/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

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

	config := config.GetConfig()

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl == "" {
		helpers.HandleJSONResponse(c, nil, errors.New("redirectUrl query parameter is required"))
		return
	}

	isUrlAllowed := false
	for _, origin := range config.Origins {
		if strings.HasPrefix(redirectUrl, origin) {
			isUrlAllowed = true
			break
		}
	}
	if !isUrlAllowed {
		helpers.HandleJSONResponse(c, nil, errors.New("redirectUrl is not allowed"))
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	url, err := ctl.signinService.GetProviderUrl(provider, redirectUrl, user)

	helpers.HandleJSONResponse(c, url, err)
}

// @Summary Get all redirect URLs for each OAuth provider
// @Tags Authentication
// @Param redirectUrl query string true "URL to redirect after OAuth authentication"
// @Success 200 {string} map[string]string "OAuth URL"
// @Failure 400 {object} helpers.ApiError
// @Router /v1/auth/providers/url [get]
func (ctl *ProviderController) ProvidersUrl(c *gin.Context) {
	config := config.GetConfig()

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl == "" {
		helpers.HandleJSONResponse(c, nil, errors.New("redirectUrl query parameter is required"))
		return
	}

	isUrlAllowed := false
	for _, origin := range config.Origins {
		if strings.HasPrefix(redirectUrl, origin) {
			isUrlAllowed = true
			break
		}
	}
	if !isUrlAllowed {
		helpers.HandleJSONResponse(c, nil, errors.New("redirectUrl is not allowed"))
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	url, err := ctl.signinService.GetProviderUrls(redirectUrl, user)

	helpers.HandleJSONResponse(c, url, err)
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

	redirectUrl := state["redirectUrl"]
	if redirectUrl == "" {
		helpers.HandleJSONResponse(c, nil, errors.New("State parameter is missing redirectUrl"))
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
