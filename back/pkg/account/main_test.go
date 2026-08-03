package account

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

// TestMain wires up the two process-global singletons this package's
func TestMain(m *testing.M) {
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	// Point at the repo's real JWT keypair so signinService.GenerateTokens
	// (used by AccountService.Create/Update) works against real RSA keys.
	if os.Getenv("AUTH_PRIVATE_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PRIVATE_PEM_PATH", "../../config/jwt/private.pem")
	}
	if os.Getenv("AUTH_PUBLIC_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PUBLIC_PEM_PATH", "../../config/jwt/public.pem")
	}
	// ForgotPassword/ResetPassword encrypt/decrypt the reset token.
	if os.Getenv("ENCRYPTION_KEY") == "" {
		_ = os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
	}
	// mail.buildEmailMessage assumes Config.Email.Address contains "@" (it
	// splits on it to build a Message-ID); without this, the background
	// welcome/reset emails sent by Create/ForgotPassword/ResetPassword panic.
	if os.Getenv("EMAIL_ADDRESS") == "" {
		_ = os.Setenv("EMAIL_ADDRESS", "noreply@slotfinder.test")
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
	db.SetDBForTests(database)

	os.Exit(m.Run())
}
