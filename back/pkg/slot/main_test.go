package slot

import (
	"app/config"
	"app/db"
	model "app/db/models"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// See the detailed comment on TestMain in pkg/account/main_test.go.
func TestMain(m *testing.M) {
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("EMAIL_ADDRESS") == "" {
		_ = os.Setenv("EMAIL_ADDRESS", "noreply@slotfinder.test")
	}
	config.Init()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(
		&model.Account{}, &model.Event{}, &model.Availability{}, &model.Slot{},
		&model.AccountEvent{}, &model.AccountProvider{}, &model.RefreshToken{},
	); err != nil {
		panic(err)
	}
	if sqlDB, err := database.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	db.SetDB(database)

	os.Exit(m.Run())
}
