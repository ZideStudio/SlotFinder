package signin

import (
	"app/config"
	"app/db"
	model "app/db/models"
	"app/testutils"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// See the detailed comment on TestMain in pkg/account/main_test.go for why
// this wiring is needed and safe for tests.
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
	if err := testutils.EnsureTestJWTKeyPair(
		os.Getenv("AUTH_PRIVATE_PEM_PATH"),
		os.Getenv("AUTH_PUBLIC_PEM_PATH"),
	); err != nil {
		panic(err)
	}
	config.Init()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(&model.Account{}, &model.AccountProvider{}, &model.RefreshToken{}); err != nil {
		panic(err)
	}
	if sqlDB, err := database.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	db.SetDBForTests(database)

	os.Exit(m.Run())
}
