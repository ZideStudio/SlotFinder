package test

import (
	model "app/db/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewTestDB opens an in-memory SQLite database migrated with every model,
// for use by repository (and DB-backed service) tests.
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}

	if err := database.AutoMigrate(
		&model.Account{},
		&model.Event{},
		&model.Availability{},
		&model.Slot{},
		&model.AccountEvent{},
		&model.AccountProvider{},
		&model.RefreshToken{},
	); err != nil {
		t.Fatalf("failed to migrate in-memory sqlite db: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return database
}
