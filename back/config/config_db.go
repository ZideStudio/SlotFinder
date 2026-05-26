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
	TimeZone string `env:"DB_TIMEZONE"`
}

// Dialect returns "postgres"
func (c DbConfiguration) Dialect() string {
	return "postgres"
}

// GetPostgresConnectionInfo returns Postgres URL string
func (c DbConfiguration) GetPostgresConnectionInfo() string {
	timezone := c.TimeZone
	if timezone == "" {
		timezone = "UTC" // Default to UTC if not specified
	}

	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable TimeZone=%s",
			c.Host, c.Port, c.User, c.Name, timezone)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, timezone)
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
		TimeZone: os.Getenv("DB_TIMEZONE"),
	}
}
