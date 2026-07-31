package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// These tests exercise TestConnection/GetDB directly against the package-private
// `conn` variable using a lightweight sqlite connection instead of a real
// Postgres instance — Ping()/Close() behave the same regardless of driver.
// Init() itself (and startMigration()) are not covered here: they require a
// live Postgres DSN and panic on failure, so they're left to integration/E2E
// testing rather than unit tests.

func TestGetDB(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	conn = database
	assert.Same(t, database, GetDB())
}

func TestTestConnection_NoConnection(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	conn = nil
	assert.False(t, TestConnection())
}

func TestTestConnection_Success(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	conn = database

	assert.True(t, TestConnection())
}

func TestTestConnection_PingFails(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	sqlDB, err := database.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Close())

	conn = database
	assert.False(t, TestConnection(), "Ping should fail on an already-closed connection")
}
