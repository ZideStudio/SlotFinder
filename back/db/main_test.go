package db

import (
	"app/config"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// loadTestDotEnv loads ../.env.test if present, so `go test` can point at a
// Postgres database dedicated to tests instead of ../.env's dev database.
// Falls back to ../.env when .env.test hasn't been set up yet (ignored if
// missing too — CI sets these env vars directly instead).
func loadTestDotEnv() {
	if _, err := os.Stat("../.env.test"); err == nil {
		_ = godotenv.Load("../.env.test")
		return
	}
	_ = godotenv.Load("../.env")
}

// TestMain loads test env vars (see loadTestDotEnv) and then requires a
// reachable Postgres instance for the whole package. db.go and migration.go
// run real Postgres-only SQL (CREATE TYPE ... AS ENUM) that sqlite can't
// execute, so every test here needs a real connection; a missing/unreachable
// Postgres fails the whole package immediately instead of silently skipping.
func TestMain(m *testing.M) {
	loadTestDotEnv()

	if _, err := strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		fmt.Fprintf(os.Stderr, "DB_PORT is not set to a valid port (%q): export DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME (or set them in back/.env.test, see back/.env.test.model) for a running Postgres instance to run this package's tests\n", os.Getenv("DB_PORT"))
		os.Exit(1)
	}

	config.Init()
	c := config.GetConfig()

	preflight, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  c.Db.GetPostgresConnectionInfo(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "no reachable Postgres at %s:%d: %v\n", c.Db.Host, c.Db.Port, err)
		os.Exit(1)
	}
	if sqlDB, dbErr := preflight.DB(); dbErr == nil {
		_ = sqlDB.Close()
	}

	os.Exit(m.Run())
}

// realTestDB opens a fresh real-Postgres connection using the package's
// config, for tests that need to manipulate the package-private `conn` var
// directly.
func realTestDB() (*gorm.DB, error) {
	c := config.GetConfig()
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  c.Db.GetPostgresConnectionInfo(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
}
