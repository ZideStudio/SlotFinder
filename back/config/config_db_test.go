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

func TestGetPostgresConfig_Success(t *testing.T) {
	envs := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USER":     "user",
		"DB_PASSWORD": "secret",
		"DB_NAME":     "db",
		"DB_TIMEZONE": "Europe/Paris",
	}
	originals := setEnvs(t, envs)
	defer restoreEnvs(originals)

	c := GetPostgresConfig()

	assert.Equal(t, "localhost", c.Host)
	assert.Equal(t, 5432, c.Port)
	assert.Equal(t, "user", c.User)
	assert.Equal(t, "secret", c.Password)
	assert.Equal(t, "db", c.Name)
	assert.Equal(t, "Europe/Paris", c.TimeZone)
}

func TestGetAuthConfig_Success(t *testing.T) {
	envs := map[string]string{
		"AUTH_PUBLIC_PEM_PATH":  "/tmp/public.pem",
		"AUTH_PRIVATE_PEM_PATH": "/tmp/private.pem",
	}
	originals := setEnvs(t, envs)
	defer restoreEnvs(originals)

	c := GetAuthConfig()

	assert.Equal(t, "/tmp/public.pem", c.PublicPemPath)
	assert.Equal(t, "/tmp/private.pem", c.PrivatePemPath)
}

func TestGetEmailConfig_Success(t *testing.T) {
	envs := map[string]string{
		"EMAIL_HOST":     "smtp.example.com",
		"EMAIL_PORT":     "587",
		"EMAIL_ADDRESS":  "noreply@example.com",
		"EMAIL_PASSWORD": "secret",
	}
	originals := setEnvs(t, envs)
	defer restoreEnvs(originals)

	c := GetEmailConfig()

	assert.Equal(t, "smtp.example.com", c.Host)
	assert.Equal(t, "587", c.Port)
	assert.Equal(t, "noreply@example.com", c.Address)
	assert.Equal(t, "secret", c.Password)
}

func TestGetProviderConfig_Success(t *testing.T) {
	envs := map[string]string{
		"PROVIDER_DISCORD_CLIENT_ID":     "discord-id",
		"PROVIDER_DISCORD_CLIENT_SECRET": "discord-secret",
		"PROVIDER_DISCORD_REDIRECT_URL":  "https://example.com/discord",
		"PROVIDER_GOOGLE_CLIENT_ID":      "google-id",
		"PROVIDER_GOOGLE_CLIENT_SECRET":  "google-secret",
		"PROVIDER_GOOGLE_REDIRECT_URL":   "https://example.com/google",
		"PROVIDER_GITHUB_CLIENT_ID":      "github-id",
		"PROVIDER_GITHUB_CLIENT_SECRET":  "github-secret",
		"PROVIDER_GITHUB_REDIRECT_URL":   "https://example.com/github",
	}
	originals := setEnvs(t, envs)
	defer restoreEnvs(originals)

	c := GetProviderConfig()

	assert.Equal(t, "discord-id", c.DiscordClientId)
	assert.Equal(t, "google-id", c.GoogleClientId)
	assert.Equal(t, "github-id", c.GithubClientId)
}

func TestInitAndGetConfig_Success(t *testing.T) {
	envs := map[string]string{
		"ENV":                   "test",
		"APP_HOST":              "0.0.0.0",
		"APP_PORT":              "8080",
		"DOMAIN":                "example.com",
		"ORIGIN":                "https://example.com",
		"DB_HOST":               "localhost",
		"DB_PORT":               "5432",
		"DB_USER":               "user",
		"DB_PASSWORD":           "secret",
		"DB_NAME":               "db",
		"DB_TIMEZONE":           "UTC",
		"AUTH_PUBLIC_PEM_PATH":  "/tmp/public.pem",
		"AUTH_PRIVATE_PEM_PATH": "/tmp/private.pem",
		"EMAIL_HOST":            "smtp.example.com",
		"EMAIL_PORT":            "587",
		"EMAIL_ADDRESS":         "noreply@example.com",
		"EMAIL_PASSWORD":        "secret",
	}
	originals := setEnvs(t, envs)
	defer restoreEnvs(originals)

	c := Init()

	assert.Equal(t, "test", c.Env)
	assert.Equal(t, "0.0.0.0", c.Host)
	assert.Equal(t, "8080", c.Port)
	assert.Equal(t, "example.com", c.Domain)
	assert.Equal(t, "https://example.com", c.Origin)
	assert.Equal(t, "localhost", c.Db.Host)
	assert.Equal(t, "/tmp/public.pem", c.Auth.PublicPemPath)
	assert.Equal(t, "smtp.example.com", c.Email.Host)
	assert.Same(t, c, GetConfig())
}

func setEnvs(t *testing.T, envs map[string]string) map[string]string {
	t.Helper()
	originals := make(map[string]string, len(envs))
	for key, value := range envs {
		originals[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	return originals
}

func restoreEnvs(originals map[string]string) {
	for key, value := range originals {
		os.Setenv(key, value)
	}
}
