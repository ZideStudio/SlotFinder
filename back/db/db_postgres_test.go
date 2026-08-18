package db

import (
	"app/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_RealPostgres(t *testing.T) {
	original := conn
	t.Cleanup(func() { conn = original })

	database := Init()
	require.NotNil(t, database)

	// Running Init() again re-enters startMigration -> ensureEventStatusEnumType
	// with the type already present, exercising its "add any missing enum
	// labels" branch, plus AutoMigrate's idempotent (no-op) path.
	database2 := Init()
	require.NotNil(t, database2)
}

// TestInit_ConnectionError_Panics points at a port nothing listens on
// (leaving DB_HOST untouched, so no DNS lookup risk) to make gorm.Open fail
// fast with "connection refused", exercising Init()'s connection-failure panic.
func TestInit_ConnectionError_Panics(t *testing.T) {
	original := conn
	t.Cleanup(func() { conn = original })

	origPort := os.Getenv("DB_PORT")
	t.Cleanup(func() {
		_ = os.Setenv("DB_PORT", origPort)
		config.Init()
	})
	require.NoError(t, os.Setenv("DB_PORT", "1"))
	config.Init()

	assert.Panics(t, func() { Init() })
}
