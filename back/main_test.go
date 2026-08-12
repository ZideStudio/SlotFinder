package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureLogging_DefaultsToInfo(t *testing.T) {
	origEnv := os.Getenv("ENV")
	origLevel := zerolog.GlobalLevel()
	t.Cleanup(func() {
		_ = os.Setenv("ENV", origEnv)
		zerolog.SetGlobalLevel(origLevel)
	})
	require.NoError(t, os.Setenv("ENV", "production"))

	configureLogging()

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
}

func TestConfigureLogging_LocalEnvUsesDebug(t *testing.T) {
	origEnv := os.Getenv("ENV")
	origLevel := zerolog.GlobalLevel()
	t.Cleanup(func() {
		_ = os.Setenv("ENV", origEnv)
		zerolog.SetGlobalLevel(origLevel)
	})
	require.NoError(t, os.Setenv("ENV", "local"))

	configureLogging()

	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}

func TestLoadDotEnvIfPresent_MissingFile_NoOp(t *testing.T) {
	err := loadDotEnvIfPresent(filepath.Join(t.TempDir(), "does-not-exist.env"))
	assert.NoError(t, err)
}

func TestLoadDotEnvIfPresent_LoadsVariables(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	require.NoError(t, os.WriteFile(path, []byte("MAIN_TEST_VAR=hello\n"), 0o600))
	t.Cleanup(func() { _ = os.Unsetenv("MAIN_TEST_VAR") })

	err := loadDotEnvIfPresent(path)

	assert.NoError(t, err)
	assert.Equal(t, "hello", os.Getenv("MAIN_TEST_VAR"))
}

func TestLoadDotEnvIfPresent_UnreadableFile(t *testing.T) {
	// A directory exists at path, so os.Stat succeeds but godotenv.Load's
	// read fails.
	dirAsEnvFile := t.TempDir()

	err := loadDotEnvIfPresent(dirAsEnvFile)

	assert.Error(t, err)
}
