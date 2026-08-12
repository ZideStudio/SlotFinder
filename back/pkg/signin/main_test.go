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
	if os.Getenv("AUTH_PRIVATE_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PRIVATE_PEM_PATH", "../../config/jwt/private.pem")
	}
	if os.Getenv("AUTH_PUBLIC_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PUBLIC_PEM_PATH", "../../config/jwt/public.pem")
	}
	if err := testutils.EnsureTestJWTKeyPair(
		os.Getenv("AUTH_PRIVATE_PEM_PATH"),
		os.Getenv("AUTH_PUBLIC_PEM_PATH"),
	); err != nil {
		panic(err)
	}
	config.Init()

	os.Exit(m.Run())
}
