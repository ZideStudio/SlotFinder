package server

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// See the detailed comment on TestMain in pkg/account/main_test.go.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	testutils.EnsureTestEncryptionEnv()
	testutils.EnsureTestEmailEnv()
	config.Init()
	testutils.MustSetupTestDB()

	os.Exit(m.Run())
}
