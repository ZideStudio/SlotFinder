package account

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// TestMain wires up the two process-global singletons this package's
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
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

	os.Exit(m.Run())
}
