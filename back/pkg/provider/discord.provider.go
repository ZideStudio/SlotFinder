package provider

import (
	"app/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DiscordUserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (s *ProviderService) getDiscordUserInfo(code string) (ProviderAccount, error) {
	providerConfig := config.GetProviderConfig()

	client := resty.New()

	var token DiscordTokenResponse
	res, err := client.R().
		SetFormData(map[string]string{
			"client_id":     providerConfig.DiscordClientId,
			"client_secret": providerConfig.DiscordClientSecret,
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  providerConfig.DiscordRedirectUrl,
		}).
		SetResult(&token).
		Post("https://discord.com/api/oauth2/token")
	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Discord token: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Discord token with status %v: %s", res.StatusCode(), res)
	}

	var userInfo DiscordUserInfo
	res, err = client.R().
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetResult(&userInfo).
		Get("https://discord.com/api/users/@me")
	if err != nil {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Discord user info: %w", err)
	}
	if res.StatusCode() != 200 {
		return ProviderAccount{}, fmt.Errorf("OAUTH: failed to get Discord user with status %v info: %s", res.StatusCode(), res)
	}

	pictureUrl := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userInfo.Id, userInfo.Avatar)

	return ProviderAccount{
		Id:        userInfo.Id,
		Username:  userInfo.Username,
		AvatarUrl: &pictureUrl,
		Email:     &userInfo.Email,
	}, nil
}
