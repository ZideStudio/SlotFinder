package config

import (
	"os"
)

type ProviderConfiguration struct {
	DiscordClientId     string `env:"PROVIDER_DISCORD_CLIENT_ID"`
	DiscordClientSecret string `env:"PROVIDER_DISCORD_CLIENT_SECRET"`
	DiscordRedirectUrl  string `env:"PROVIDER_DISCORD_REDIRECT_URL"`

	GoogleClientId     string `env:"PROVIDER_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"PROVIDER_GOOGLE_CLIENT_SECRET"`
	GoogleRedirectUrl  string `env:"PROVIDER_GOOGLE_REDIRECT_URL"`

	GithubClientId     string `env:"PROVIDER_GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"PROVIDER_GITHUB_CLIENT_SECRET"`
	GithubRedirectUrl  string `env:"PROVIDER_GITHUB_REDIRECT_URL"`
}

func GetProviderConfig() ProviderConfiguration {
	return ProviderConfiguration{
		DiscordClientId:     os.Getenv("PROVIDER_DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("PROVIDER_DISCORD_CLIENT_SECRET"),
		DiscordRedirectUrl:  os.Getenv("PROVIDER_DISCORD_REDIRECT_URL"),
		GoogleClientId:      os.Getenv("PROVIDER_GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("PROVIDER_GOOGLE_CLIENT_SECRET"),
		GoogleRedirectUrl:   os.Getenv("PROVIDER_GOOGLE_REDIRECT_URL"),
		GithubClientId:      os.Getenv("PROVIDER_GITHUB_CLIENT_ID"),
		GithubClientSecret:  os.Getenv("PROVIDER_GITHUB_CLIENT_SECRET"),
		GithubRedirectUrl:   os.Getenv("PROVIDER_GITHUB_REDIRECT_URL"),
	}
}
