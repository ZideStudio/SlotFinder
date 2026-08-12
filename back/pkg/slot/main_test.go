package slot

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// See the detailed comment on TestMain in pkg/account/main_test.go.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	if os.Getenv("EMAIL_ADDRESS") == "" {
		_ = os.Setenv("EMAIL_ADDRESS", "noreply@slotfinder.test")
	}
	config.Init()

	os.Exit(m.Run())
}
