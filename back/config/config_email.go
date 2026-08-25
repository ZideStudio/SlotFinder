package config

import (
	"fmt"
	"os"
	"strconv"
)

type EmailConfiguration struct {
	Host     string `env:"EMAIL_HOST"`
	Port     string `env:"EMAIL_PORT"`
	Address  string `env:"EMAIL_ADDRESS"`
	Password string `env:"EMAIL_PASSWORD"`
	Disabled bool   `env:"EMAIL_SENDING_DISABLED"`
}

func GetEmailConfig() EmailConfiguration {
	return EmailConfiguration{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Address:  os.Getenv("EMAIL_ADDRESS"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Disabled: getEmailSendingDisabled(),
	}
}

func getEmailSendingDisabled() bool {
	raw := os.Getenv("EMAIL_SENDING_DISABLED")
	if raw == "" {
		return false
	}

	disabled, err := strconv.ParseBool(raw)
	if err != nil {
		panic(fmt.Errorf("invalid EMAIL_SENDING_DISABLED value %q: %w", raw, err))
	}

	return disabled
}
