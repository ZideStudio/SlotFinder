package guard

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	testutils.EnsureTestAuthEnv()
	if os.Getenv("ORIGIN") == "" {
		_ = os.Setenv("ORIGIN", "https://slotfinder.test")
	}
	config.Init()

	os.Exit(m.Run())
}
