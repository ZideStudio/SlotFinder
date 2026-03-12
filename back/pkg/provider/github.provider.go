package provider

import (
	"app/config"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GithubUserNameInfo struct {
	Id        int32  `json:"id"`
	Name      string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

type GithubUserEmailInfo struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *ProviderService) getGithubUserInfo(code string) (ProviderAccount, error) {
	providerConfig := config.GetProviderConfig()

	client := resty.New()

	var token GithubTokenResponse
	res, err := client.R().
		SetFormData(map[string]string{
			"client_id":     providerConfig.GithubClientId,
			"client_secret": providerConfig.GithubClientSecret,
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  providerConfig.GithubRedirectUrl,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&token).
		Post("https://github.com/login/oauth/access_token")
	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github token: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github token with status %v: %s", res.StatusCode(), res)
	}

	var userNameInfo GithubUserNameInfo
	res, err = client.R().
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetResult(&userNameInfo).
		Get("https://api.github.com/user")

	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github user info: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github user with status %v info: %s", res.StatusCode(), res)
	}

	var emailInfo []GithubUserEmailInfo
	res, err = client.R().
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetResult(&emailInfo).
		Get("https://api.github.com/user/emails")

	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github emails info: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Github emails with status %v info: %s", res.StatusCode(), res)
	}

	if len(emailInfo) == 0 {
		return ProviderAccount{}, errors.New("OAUTH: no email addresses found for Github user")
	}
	var email string
	for _, currentEmail := range emailInfo {
		if currentEmail.Primary {
			email = currentEmail.Email
			break
		}
	}
	if email == "" {
		email = emailInfo[0].Email
	}

	return ProviderAccount{
		Id:        fmt.Sprintf("%d", userNameInfo.Id),
		Username:  userNameInfo.Name,
		AvatarUrl: &userNameInfo.AvatarUrl,
		Email:     &email,
	}, nil
}
