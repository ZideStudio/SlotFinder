package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbConfiguration_Dialect(t *testing.T) {
	c := DbConfiguration{}
	assert.Equal(t, "postgres", c.Dialect())
}

func TestDbConfiguration_GetPostgresConnectionInfo_DefaultsTimeZoneToUTC(t *testing.T) {
	c := DbConfiguration{Host: "localhost", Port: 5432, User: "user", Name: "db"}
	info := c.GetPostgresConnectionInfo()
	assert.Contains(t, info, "TimeZone=UTC")
	assert.NotContains(t, info, "password=")
}

func TestDbConfiguration_GetPostgresConnectionInfo_WithPassword(t *testing.T) {
	c := DbConfiguration{Host: "localhost", Port: 5432, User: "user", Password: "secret", Name: "db", TimeZone: "Europe/Paris"}
	info := c.GetPostgresConnectionInfo()
	assert.Contains(t, info, "password=secret")
	assert.Contains(t, info, "TimeZone=Europe/Paris")
}

func TestGetPostgresConfig_InvalidPort_Panics(t *testing.T) {
	original := os.Getenv("DB_PORT")
	defer os.Setenv("DB_PORT", original)

	os.Setenv("DB_PORT", "not-a-number")

	assert.Panics(t, func() {
		GetPostgresConfig()
	})
}
