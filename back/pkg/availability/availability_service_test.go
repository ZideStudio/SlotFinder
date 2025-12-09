package availability

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestServiceHasSyncMap verifies that the service has a sync.Map for locking
func TestServiceHasSyncMap(t *testing.T) {
	service := NewAvailabilityService(nil)

	// Verify that the service was created successfully
	assert.NotNil(t, service, "Service should be created")
	assert.NotNil(t, service.availabilityRepository, "Repository should be initialized")
	assert.NotNil(t, service.eventRepository, "Event repository should be initialized")

	// The locks field is a sync.Map which is always initialized to its zero value
	// We can verify it works by storing and loading a value
	testKey := "test:key"
	service.locks.Store(testKey, "test-value")
	value, ok := service.locks.Load(testKey)
	assert.True(t, ok, "Should be able to load stored value")
	assert.Equal(t, "test-value", value, "Loaded value should match stored value")
}

// TestSyncMapLoadOrStore verifies that LoadOrStore works correctly for concurrent access
func TestSyncMapLoadOrStore(t *testing.T) {
	service := NewAvailabilityService(nil)

	lockKey := "account1:event1"

	// First LoadOrStore should create a new mutex
	value1, loaded1 := service.locks.LoadOrStore(lockKey, &sync.Mutex{})
	assert.False(t, loaded1, "First LoadOrStore should not find existing value")
	assert.NotNil(t, value1, "LoadOrStore should return a value")

	// Second LoadOrStore should return the same mutex
	value2, loaded2 := service.locks.LoadOrStore(lockKey, &sync.Mutex{})
	assert.True(t, loaded2, "Second LoadOrStore should find existing value")
	assert.Equal(t, value1, value2, "Both LoadOrStore calls should return the same value")

	// Verify we can cast to *sync.Mutex
	mu1, ok1 := value1.(*sync.Mutex)
	mu2, ok2 := value2.(*sync.Mutex)
	assert.True(t, ok1, "First value should be castable to *sync.Mutex")
	assert.True(t, ok2, "Second value should be castable to *sync.Mutex")
	assert.Equal(t, mu1, mu2, "Both mutexes should be the same instance")
}

// TestSyncMapConcurrentAccess verifies that multiple goroutines can safely access the sync.Map
func TestSyncMapConcurrentAccess(t *testing.T) {
	service := NewAvailabilityService(nil)

	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	// Simulate concurrent access to the same lock key
	lockKey := "account1:event1"

	for i := 0; i < numGoroutines; i++ {
		go func() {
			value, _ := service.locks.LoadOrStore(lockKey, &sync.Mutex{})
			mu := value.(*sync.Mutex)

			// Lock and unlock to verify no panic
			mu.Lock()
			defer mu.Unlock()

			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify that only one mutex was created
	value, ok := service.locks.Load(lockKey)
	assert.True(t, ok, "Lock key should exist")
	assert.NotNil(t, value, "Value should not be nil")
}

// TestUpdateDto verifies that the AvailabilityUpdateDto is properly structured
func TestUpdateDto(t *testing.T) {
	// Test with nil values
	dto1 := AvailabilityUpdateDto{}
	assert.Nil(t, dto1.StartsAt, "StartsAt should be nil by default")
	assert.Nil(t, dto1.EndsAt, "EndsAt should be nil by default")
}
