package guard

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	if os.Getenv("ORIGIN") == "" {
		_ = os.Setenv("ORIGIN", "https://slotfinder.test")
	}
	config.Init()

	os.Exit(m.Run())
}
