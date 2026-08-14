package signin

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// See the detailed comment on TestMain in pkg/account/main_test.go for why
// this wiring is needed and safe for tests.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	config.Init()

	os.Exit(m.Run())
}
