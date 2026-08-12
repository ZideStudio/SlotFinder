package testutils

import (
	"app/config"
	"app/db"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// migrationLockKey is an arbitrary Postgres advisory lock id, unique to this
// project's test schema migration.
const migrationLockKey = 868012025

var (
	baseOnce sync.Once
	baseDB   *gorm.DB
	baseErr  error
)

// LoadTestEnv loads .env.test (or .env) from the module root regardless of
// which package's directory `go test` runs from, and defaults DB_PORT so
// config.Init doesn't panic when neither file exists (CI sets DB_* directly).
func LoadTestEnv() {
	if dir := findModuleRoot(); dir != "" {
		envTest := filepath.Join(dir, ".env.test")
		if _, err := os.Stat(envTest); err == nil {
			_ = godotenv.Load(envTest)
		} else {
			_ = godotenv.Load(filepath.Join(dir, ".env"))
		}
	}
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
}

func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// ensureBase connects+migrates the shared test Postgres database once per
// binary. The advisory lock serializes that migration DDL across the
// separate `go test` binaries each package runs as, all hitting one database.
func ensureBase() (*gorm.DB, error) {
	baseOnce.Do(func() {
		LoadTestEnv()
		config.Init()

		c := config.GetConfig()
		lockConn, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  c.Db.GetPostgresConnectionInfo(),
			PreferSimpleProtocol: true,
		}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			baseErr = err
			return
		}
		if err := lockConn.Exec("SELECT pg_advisory_lock(?)", migrationLockKey).Error; err != nil {
			baseErr = err
			return
		}
		defer func() {
			lockConn.Exec("SELECT pg_advisory_unlock(?)", migrationLockKey)
			if sqlDB, err := lockConn.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}()

		baseDB = db.Init()
	})
	return baseDB, baseErr
}

func requireBase(t *testing.T) *gorm.DB {
	t.Helper()
	base, err := ensureBase()
	if base == nil {
		t.Fatalf("test postgres database unavailable (see back/README.md#test): %v", err)
	}
	return base
}

// TestDB begins a transaction on its own dedicated connection (never
// reused by another test — see openDedicatedConn), rolled back and closed
// via t.Cleanup, and points db.GetDB() at it for the test's duration.
func TestDB(t *testing.T) *gorm.DB {
	t.Helper()
	base := requireBase(t)

	conn := openDedicatedConn(t)
	tx := conn.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin test transaction: %v", tx.Error)
	}
	db.SetDBForTests(tx)
	t.Cleanup(func() {
		tx.Rollback()
		if sqlDB, err := conn.DB(); err == nil {
			_ = sqlDB.Close()
		}
		db.SetDBForTests(base)
	})
	return tx
}

// FreshDB opens a connection independent of any transaction — for tests
// that close it themselves (e.g. exercising db.TestConnection(), which
// closes whatever pool it's given).
func FreshDB(t *testing.T) *gorm.DB {
	t.Helper()
	requireBase(t)
	return openDedicatedConn(t)
}

// openDedicatedConn pins a single physical connection so it's never handed
// to another test: a straggling goroutine from an async call (see
// AwaitAsyncDBWork) could otherwise corrupt a connection reused later.
func openDedicatedConn(t *testing.T) *gorm.DB {
	t.Helper()
	c := config.GetConfig()
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  c.Db.GetPostgresConnectionInfo(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("failed to open postgres connection: %v", err)
	}
	if sqlDB, err := database.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	return database
}

// ClosedDB returns a *gorm.DB whose connection is already closed, so any
// query fails immediately — forces one repository's calls to fail while
// the rest of a service keeps working.
func ClosedDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := FreshDB(t)
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}
	return database
}

// BaseDB returns the shared, non-transactional base connection to the test
// Postgres database — for restoring db.GetDB() after a test points it at its
// own connection (see FreshDB).
func BaseDB(t *testing.T) *gorm.DB {
	t.Helper()
	return requireBase(t)
}

// MustSetupTestDB is TestMain's entry point (no *testing.T there): same as
// ensureBase but panics on failure, for packages that build routers/services
// once per test instead of a fresh transaction per test.
func MustSetupTestDB() *gorm.DB {
	base, err := ensureBase()
	if base == nil {
		panic("test postgres database unavailable (see back/README.md#test): " + err.Error())
	}
	return base
}

// MakeReadOnly flips db.GetDB() read-only so a subsequent write fails while
// prior reads keep working; undone automatically on transaction rollback.
func MakeReadOnly(t *testing.T) {
	t.Helper()
	if err := db.GetDB().Exec("SET transaction_read_only = on").Error; err != nil {
		t.Fatalf("failed to set transaction_read_only: %v", err)
	}
}

// AwaitAsyncDBWork gives a fire-and-forget goroutine (e.g. "go
// s.slotService.LoadSlots(...)") time to finish before the test queries
// again — it shares this test's connection, and pgx isn't concurrency-safe.
func AwaitAsyncDBWork(t *testing.T) {
	t.Helper()
	time.Sleep(100 * time.Millisecond)
}
