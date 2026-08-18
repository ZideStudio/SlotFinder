package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Exercises TestConnection/GetDB against the package-private `conn` var
// using a real Postgres connection (see TestMain).

func TestGetDB(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	database, err := realTestDB()
	require.NoError(t, err)

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

	database, err := realTestDB()
	require.NoError(t, err)
	conn = database

	assert.True(t, TestConnection())
}

func TestTestConnection_DBFails(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	// A gorm.DB with no ConnPool set can't be converted to a *sql.DB.
	conn = &gorm.DB{Config: &gorm.Config{}}
	assert.False(t, TestConnection())
}

func TestTestConnection_PingFails(t *testing.T) {
	original := conn
	defer func() { conn = original }()

	database, err := realTestDB()
	require.NoError(t, err)

	sqlDB, err := database.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	conn = database
	assert.False(t, TestConnection(), "Ping should fail on an already-closed connection")
}
