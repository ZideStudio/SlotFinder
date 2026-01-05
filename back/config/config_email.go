package config

import (
	"os"
)

type EmailConfiguration struct {
	Host     string `env:"EMAIL_HOST"`
	Port     string `env:"EMAIL_PORT"`
	Address  string `env:"EMAIL_ADDRESS"`
	Password string `env:"EMAIL_PASSWORD"`
}

func GetEmailConfig() EmailConfiguration {
	return EmailConfiguration{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Address:  os.Getenv("EMAIL_ADDRESS"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}
