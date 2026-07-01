package config

import (
	"os"
)

type Config struct {
	Env      string `env:"ENV"`
	Host     string `env:"APP_HOST"`
	Port     string `env:"APP_PORT"`
	Domain   string `env:"DOMAIN"`
	Origin   string `env:"ORIGIN"`
	Db       DbConfiguration
	Auth     AuthConfiguration
	Provider ProviderConfiguration
	Email    EmailConfiguration
}

var config *Config

func Init() *Config {
	config = &Config{
		Env:      os.Getenv("ENV"),
		Db:       GetPostgresConfig(),
		Host:     os.Getenv("APP_HOST"),
		Port:     os.Getenv("APP_PORT"),
		Domain:   os.Getenv("DOMAIN"),
		Origin:   os.Getenv("ORIGIN"),
		Auth:     GetAuthConfig(),
		Provider: GetProviderConfig(),
		Email:    GetEmailConfig(),
	}

	return config
}

func GetConfig() *Config {
	return config
}
