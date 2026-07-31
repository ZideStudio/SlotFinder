package provider

import (
	"app/config"
	"app/db"
	model "app/db/models"
	"app/testutils"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestMain wires the process-global config/DB singletons this package's
// NewXService(nil) fallbacks depend on (ProviderService composes
// signinService, accountService, avatarService, mailService — all built the
// same way). See the equivalent, more detailed comment in
// pkg/account/main_test.go for why this is necessary and safe for tests.
func TestMain(m *testing.M) {
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("AUTH_PRIVATE_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PRIVATE_PEM_PATH", "../../config/jwt/private.pem")
	}
	if os.Getenv("AUTH_PUBLIC_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PUBLIC_PEM_PATH", "../../config/jwt/public.pem")
	}
	if os.Getenv("ENCRYPTION_KEY") == "" {
		_ = os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
	}
	if os.Getenv("EMAIL_ADDRESS") == "" {
		_ = os.Setenv("EMAIL_ADDRESS", "noreply@slotfinder.test")
	}
	if os.Getenv("PROVIDER_GOOGLE_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_GOOGLE_REDIRECT_URL", "https://slotfinder.test/oauth/google/callback")
	}
	if os.Getenv("PROVIDER_DISCORD_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_DISCORD_REDIRECT_URL", "https://slotfinder.test/oauth/discord/callback")
	}
	if os.Getenv("PROVIDER_GITHUB_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_GITHUB_REDIRECT_URL", "https://slotfinder.test/oauth/github/callback")
	}
	if err := testutils.EnsureTestJWTKeyPair(
		os.Getenv("AUTH_PRIVATE_PEM_PATH"),
		os.Getenv("AUTH_PUBLIC_PEM_PATH"),
	); err != nil {
		panic(err)
	}
	config.Init()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(&model.Account{}, &model.AccountProvider{}, &model.RefreshToken{}); err != nil {
		panic(err)
	}
	if sqlDB, err := database.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	db.SetDB(database)

	os.Exit(m.Run())
}
