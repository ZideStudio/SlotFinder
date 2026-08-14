package testutils

import (
	"app/config"
	"app/db"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fatalCalled is the sentinel panic fakeT.Fatalf raises instead of halting
// via runtime.Goexit like the real testing.T.Fatalf does.
type fatalCalled struct{ message string }

// fakeT is a minimal testingT double, letting tests drive Fatalf branches
// without failing the real test binary.
type fakeT struct{}

func (fakeT) Helper()        {}
func (fakeT) Cleanup(func()) {}
func (fakeT) Fatalf(format string, args ...any) {
	panic(fatalCalled{message: fmt.Sprintf(format, args...)})
}

// expectFatal runs fn and returns its Fatalf message; a non-sentinel panic
// is re-raised as a genuine bug.
func expectFatal(t *testing.T, fn func()) (message string) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected Fatalf to be called, but the function returned normally")
			return
		}
		fc, ok := r.(fatalCalled)
		if !ok {
			panic(r)
		}
		message = fc.message
	}()
	fn()
	return
}

// resetBaseState clears the shared ensureBase() memoization, restoring it
// after the test so later tests get a fresh connection.
func resetBaseState(t *testing.T) {
	t.Helper()
	baseOnce = sync.Once{}
	baseDB = nil
	baseErr = nil
	t.Cleanup(func() {
		baseOnce = sync.Once{}
		baseDB = nil
		baseErr = nil
	})
}

func TestLoadTestEnv_FallsBackToDotEnvAndDefaultsDBPort(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module spike\n\ngo 1.26\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("PGTEST_ENV_FALLBACK_MARKER=loaded\n"), 0o600))

	require.NoError(t, os.Unsetenv("PGTEST_ENV_FALLBACK_MARKER"))
	t.Cleanup(func() { _ = os.Unsetenv("PGTEST_ENV_FALLBACK_MARKER") })

	origPort, hadPort := os.LookupEnv("DB_PORT")
	require.NoError(t, os.Unsetenv("DB_PORT"))
	t.Cleanup(func() {
		if hadPort {
			_ = os.Setenv("DB_PORT", origPort)
		} else {
			_ = os.Unsetenv("DB_PORT")
		}
	})

	t.Chdir(dir)

	LoadTestEnv()

	assert.Equal(t, "loaded", os.Getenv("PGTEST_ENV_FALLBACK_MARKER"), "should fall back to .env when .env.test is absent")
	assert.Equal(t, "5432", os.Getenv("DB_PORT"), "should default DB_PORT when unset")
}

func TestFindModuleRoot_ReturnsEmptyWhenNoGoModAncestor(t *testing.T) {
	t.Chdir(t.TempDir())

	assert.Equal(t, "", findModuleRoot())
}

func TestEnsureBase_ReturnsErrorWhenDatabaseUnreachable(t *testing.T) {
	resetBaseState(t)
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "1")

	base, err := ensureBase()

	assert.Nil(t, base)
	assert.Error(t, err)
}

func TestMustSetupTestDB_PanicsWhenDatabaseUnreachable(t *testing.T) {
	resetBaseState(t)
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "1")

	assert.Panics(t, func() {
		MustSetupTestDB()
	})
}

func TestRequireBase_FatalsWhenDatabaseUnreachable(t *testing.T) {
	resetBaseState(t)
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "1")

	msg := expectFatal(t, func() {
		requireBase(fakeT{})
	})

	assert.Contains(t, msg, "test postgres database unavailable")
}

func TestOpenDedicatedConn_FatalsWhenDatabaseUnreachable(t *testing.T) {
	requireBase(t) // establish a real base connection before breaking config below

	cfg := config.GetConfig()
	orig := cfg.Db
	cfg.Db = config.DbConfiguration{
		Host: "127.0.0.1", Port: 1,
		User: orig.User, Password: orig.Password, Name: orig.Name, TimeZone: orig.TimeZone,
	}
	t.Cleanup(func() { cfg.Db = orig })

	msg := expectFatal(t, func() {
		openDedicatedConn(fakeT{})
	})

	assert.Contains(t, msg, "failed to open postgres connection")
}

func TestMakeReadOnly_FatalsWhenConnectionClosed(t *testing.T) {
	base := requireBase(t)
	closed := ClosedDB(t)

	db.SetDBForTests(closed)
	t.Cleanup(func() { db.SetDBForTests(base) })

	msg := expectFatal(t, func() {
		MakeReadOnly(fakeT{})
	})

	assert.Contains(t, msg, "failed to set transaction_read_only")
}
