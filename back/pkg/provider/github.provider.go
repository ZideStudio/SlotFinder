package provider

import (
	"app/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GithubUserInfo struct {
	Id    string
	Name  string
	Email string
}

type GithubUserNameInfo struct {
	Id    int32  `json:"id"`
	Name  string `json:"login"`
	Email string `json:"email"`
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

	return ProviderAccount{
		Id:       fmt.Sprintf("%d", userNameInfo.Id),
		Username: userNameInfo.Name,
		Email:    userNameInfo.Email,
	}, nil
}
