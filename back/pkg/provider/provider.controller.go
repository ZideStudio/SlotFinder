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
	"net/http"
	"net/url"

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
// @Param returnUrl query string false "URL to return to after OAuth"
// @Success 200 {string} string "OAuth URL"
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/auth/{provider}/url [get]
func (ctl *ProviderController) ProviderUrl(c *gin.Context) {
	provider := c.Param("provider")
	returnUrl := c.Query("returnUrl")
	if returnUrl != "" && returnUrl[0] != '/' {
		helpers.HandleJSONResponse(c, nil, errors.New("returnUrl must be a relative path starting with /"))
		return
	}

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
		helpers.HandleJSONResponse(c, nil, errors.New("state parameter is missing"))
		return
	}
	// decrypt state
	decodedState, err := encryption.Decrypt(stateStrEncrypted)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, errors.New("state parameter has invalid encryption"))
		return
	}

	var state map[string]string
	if err := json.Unmarshal([]byte(decodedState), &state); err != nil {
		helpers.HandleJSONResponse(c, nil, errors.New("state parameter contains invalid JSON"))
		return
	} else if len(state) == 0 {
		helpers.HandleJSONResponse(c, nil, errors.New("state parameter is invalid"))
		return
	}

	redirectUrl := config.GetConfig().Origin + "/oauth/callback"
	userId := state["userId"]
	returnUrl := state["returnUrl"]
	q := url.Values{}
	q.Set("returnUrl", returnUrl)

	jwt, err := ctl.signinService.ProviderCallback(provider, code, userId)
	if err != nil {
		q.Set("error", constants.ERR_PROVIDER_CONNECTION_FAILED.Err.Error())
		redirectWithQuery := redirectUrl + "?" + q.Encode()
		c.Redirect(http.StatusFound, redirectWithQuery)
		return
	}

	lib.SetAccessTokenCookie(c, jwt.AccessToken, 0)
	lib.SetRefreshTokenCookie(c, jwt.RefreshToken, 0)

	redirectWithQuery := redirectUrl + "?" + q.Encode()
	c.Redirect(http.StatusFound, redirectWithQuery)
}
