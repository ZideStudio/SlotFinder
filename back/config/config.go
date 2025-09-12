package config

import (
	"os"
)

type Config struct {
	Env      string `env:"ENV"`
	Host     string `env:"APP_HOST"`
	Port     string `env:"APP_PORT"`
	Db       DbConfiguration
	Auth     AuthConfiguration
	Provider ProviderConfiguration
}

var config *Config

func Init() *Config {
	config = &Config{
		Env:      os.Getenv("ENV"),
		Db:       GetPostgresConfig(),
		Host:     os.Getenv("APP_HOST"),
		Port:     os.Getenv("APP_PORT"),
		Auth:     GetAuthConfig(),
		Provider: GetProviderConfig(),
	}

	return config
}

func GetConfig() *Config {
	return config
}
