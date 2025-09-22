package config

import (
	"os"
)

type AuthConfiguration struct {
	PublicPemPath  string `env:"AUTH_PUBLIC_PEM_PATH"`
	PrivatePemPath string `env:"AUTH_PRIVATE_PEM_PATH"`
}

func GetAuthConfig() AuthConfiguration {

	return AuthConfiguration{
		PublicPemPath:  os.Getenv("AUTH_PUBLIC_PEM_PATH"),
		PrivatePemPath: os.Getenv("AUTH_PRIVATE_PEM_PATH"),
	}
}
