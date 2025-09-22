package config

import (
	"fmt"
	"os"
	"strconv"
)

type DbConfiguration struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

// Dialect returns "postgres"
func (c DbConfiguration) Dialect() string {
	return "postgres"
}

// GetPostgresConnectionInfo returns Postgres URL string
func (c DbConfiguration) GetPostgresConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, "Europe/Paris")
}

// GetPostgresConfig returns PostgresConfig object
func GetPostgresConfig() DbConfiguration {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return DbConfiguration{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
