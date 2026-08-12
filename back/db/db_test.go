package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests exercise TestConnection/GetDB directly against the
// package-private `conn` variable using a real Postgres connection (see
// TestMain) — Ping()/Close() are driver behavior, but Init() itself requires
// the real thing, so everything in this package runs against Postgres now.

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
