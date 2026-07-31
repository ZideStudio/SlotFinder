package guard

import (
	"app/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("AUTH_PRIVATE_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PRIVATE_PEM_PATH", "../../config/jwt/private.pem")
	}
	if os.Getenv("AUTH_PUBLIC_PEM_PATH") == "" {
		_ = os.Setenv("AUTH_PUBLIC_PEM_PATH", "../../config/jwt/public.pem")
	}
	if os.Getenv("ORIGIN") == "" {
		_ = os.Setenv("ORIGIN", "https://slotfinder.test")
	}
	config.Init()

	os.Exit(m.Run())
}
