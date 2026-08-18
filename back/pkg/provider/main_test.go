package provider

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// TestMain wires the config this package's NewXService(nil) fallbacks
// depend on (see pkg/account/main_test.go for the detailed rationale).
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	testutils.EnsureTestEncryptionEnv()
	testutils.EnsureTestEmailEnv()
	if os.Getenv("PROVIDER_GOOGLE_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_GOOGLE_REDIRECT_URL", "https://slotfinder.test/oauth/google/callback")
	}
	if os.Getenv("PROVIDER_DISCORD_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_DISCORD_REDIRECT_URL", "https://slotfinder.test/oauth/discord/callback")
	}
	if os.Getenv("PROVIDER_GITHUB_REDIRECT_URL") == "" {
		_ = os.Setenv("PROVIDER_GITHUB_REDIRECT_URL", "https://slotfinder.test/oauth/github/callback")
	}
	config.Init()

	os.Exit(m.Run())
}
