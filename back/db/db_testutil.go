package db

import (
	"testing"

	"gorm.io/gorm"
)

// SetDBForTests overrides the process-global connection. It exists so tests
// in other packages can point every `repository.NewXRepository(nil)` /
// `NewXService(nil)` fallback at an in-memory test database instead of a
// real Postgres instance, without needing to thread a *gorm.DB through
// every nested service by hand.
//
// Panics if called outside of a test binary.
func SetDBForTests(database *gorm.DB) {
	if !testing.Testing() {
		panic("db.SetDBForTests must only be called from tests")
	}
	conn = database
}
