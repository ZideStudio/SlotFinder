package provider

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// TestMain wires the process-global config singleton this package's
// NewXService(nil) fallbacks depend on (ProviderService composes
// signinService, accountService, avatarService, mailService — all built the
// same way). See the equivalent, more detailed comment in
// pkg/account/main_test.go for why this is necessary and safe for tests.
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
