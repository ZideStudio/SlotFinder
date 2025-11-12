package config

import (
	"os"
	"strings"
)

type Config struct {
	Env         string `env:"ENV"`
	Host        string `env:"APP_HOST"`
	Port        string `env:"APP_PORT"`
	Domain      string `env:"DOMAIN"`
	ImgBBApiKey string `env:"IMGBB_API_KEY"`
	Origins     []string
	Db          DbConfiguration
	Auth        AuthConfiguration
	Provider    ProviderConfiguration
}

var config *Config

func Init() *Config {
	config = &Config{
		Env:         os.Getenv("ENV"),
		Db:          GetPostgresConfig(),
		Host:        os.Getenv("APP_HOST"),
		Port:        os.Getenv("APP_PORT"),
		Domain:      os.Getenv("DOMAIN"),
		ImgBBApiKey: os.Getenv("IMGBB_API_KEY"),
		Origins:     GetOrigin(),
		Auth:        GetAuthConfig(),
		Provider:    GetProviderConfig(),
	}

	return config
}

func GetConfig() *Config {
	return config
}

func GetOrigin() (origins []string) {
	originsString := os.Getenv("ORIGINS")
	if originsString != "" {
		origins = strings.Split(originsString, ",")
	}

	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	if len(origins) == 0 {
		panic("ORIGINS environment variable is not set")
	}

	return origins
}
