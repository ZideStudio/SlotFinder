package account

import (
	"app/config"
	"app/testutils"
	"os"
	"testing"
)

// TestMain wires up the process-global singletons this package's
// AccountService.Create/Update/ForgotPassword/ResetPassword depend on: a
// real JWT keypair (signinService.GenerateTokens), an encryption key
// (reset-token encrypt/decrypt), and an email address (mail.buildEmailMessage
// splits it on "@" to build a Message-ID) for the background emails they send.
func TestMain(m *testing.M) {
	testutils.LoadTestEnv()
	testutils.EnsureTestAuthEnv()
	testutils.EnsureTestEncryptionEnv()
	testutils.EnsureTestEmailEnv()
	config.Init()

	os.Exit(m.Run())
}
