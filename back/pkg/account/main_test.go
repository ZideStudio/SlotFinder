package account

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// TestMain wires up a real JWT keypair, encryption key, and email address —
// all of which AccountService's tests depend on.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	testutils.EnsureTestEncryptionEnv()
	testutils.EnsureTestEmailEnv()
	config.Init()

	os.Exit(m.Run())
}
