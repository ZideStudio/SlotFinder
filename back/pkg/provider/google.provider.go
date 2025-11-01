package provider

import (
	"app/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GoogleUserInfo struct {
	Id       string `json:"sub"`
	Username string `json:"name"`
	Email    string `json:"email"`
}

type GoogleTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *ProviderService) getGoogleUserInfo(code string) (ProviderAccount, error) {
	providerConfig := config.GetProviderConfig()

	client := resty.New()

	var token GoogleTokenResponse
	res, err := client.R().
		SetFormData(map[string]string{
			"client_id":     providerConfig.GoogleClientId,
			"client_secret": providerConfig.GoogleClientSecret,
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  providerConfig.GoogleRedirectUrl,
		}).
		SetResult(&token).
		Post("https://oauth2.googleapis.com/token")
	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Google token: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Google token with status %v: %s", res.StatusCode(), res)
	}

	var userInfo GoogleUserInfo
	res, err = client.R().
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetResult(&userInfo).
		Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Google user info: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Google user with status %v info: %s", res.StatusCode(), res)
	}

	return ProviderAccount(userInfo), nil
}
