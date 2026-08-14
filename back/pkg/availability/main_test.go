package availability

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// See the detailed comment on TestMain in pkg/account/main_test.go.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestEmailEnv()
	config.Init()

	os.Exit(m.Run())
}
