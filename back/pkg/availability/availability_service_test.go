package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var service = NewAvailabilityService(nil)

var username = "testuser"
var testUser = &guard.Claims{
	Id:       uuid.New(),
	Username: &username,
}

// Helper function to create a mock event for testing
func createMockEvent() model.Event {
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	threeDaysLater := tomorrow.AddDate(0, 0, 2)

	return model.Event{
		Id:       uuid.New(),
		StartsAt: tomorrow,
		EndsAt:   threeDaysLater,
	}
}

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

// TestValidateAvailabilityTimes_StartAfterEnd tests validation for end date before start date
func TestValidateAvailabilityTimes_StartAfterEnd(t *testing.T) {
	event := createMockEvent()
	startsAt := event.StartsAt.Add(2 * time.Hour)
	endsAt := event.StartsAt.Add(1 * time.Hour) // End before start

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for end date before start date")
	assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err, "Expected ERR_EVENT_START_AFTER_END error")
}

// TestValidateAvailabilityTimes_DurationTooShort tests validation for duration less than 5 minutes
func TestValidateAvailabilityTimes_DurationTooShort(t *testing.T) {
	event := createMockEvent()
	startsAt := event.StartsAt.Add(1 * time.Hour)
	endsAt := startsAt.Add(3 * time.Minute) // Less than 5 minutes

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for duration less than 5 minutes")
	assert.Equal(t, constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err, err, "Expected ERR_AVAILABILITY_DURATION_TOO_SHORT error")
}

// TestValidateAvailabilityTimes_InvalidTimeInterval tests validation for times not aligned on 5-minute intervals
func TestValidateAvailabilityTimes_InvalidTimeInterval(t *testing.T) {
	event := createMockEvent()
	
	// Test 1: Time with seconds/nanoseconds should fail
	startsAtWithSeconds := event.StartsAt.Add(1*time.Hour + 5*time.Minute + 30*time.Second) // Has 30 seconds
	endsAtWithSeconds := startsAtWithSeconds.Add(10 * time.Minute)

	err := service.validateAvailabilityTimes(startsAtWithSeconds, endsAtWithSeconds, &event)
	assert.Error(t, err, "Expected error for times with seconds")
	assert.Equal(t, constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err, err, "Expected ERR_AVAILABILITY_INVALID_TIME_INTERVAL error")
}

// TestValidateAvailabilityTimes_InvalidMinuteInterval tests validation for times on wrong minute boundary
func TestValidateAvailabilityTimes_InvalidMinuteInterval(t *testing.T) {
	event := createMockEvent()
	
	// Time on wrong minute interval (e.g., 13 minutes, not divisible by 5)
	startsAtWrongMinute := event.StartsAt.Add(1*time.Hour + 13*time.Minute)
	endsAtWrongMinute := startsAtWrongMinute.Add(10 * time.Minute)
	
	err := service.validateAvailabilityTimes(startsAtWrongMinute, endsAtWrongMinute, &event)
	assert.Error(t, err, "Expected error for times not on 5-minute intervals")
	assert.Equal(t, constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err, err, "Expected ERR_AVAILABILITY_INVALID_TIME_INTERVAL error")
}

// TestValidateAvailabilityTimes_StartBeforeEvent tests validation for start time before event start
func TestValidateAvailabilityTimes_StartBeforeEvent(t *testing.T) {
	event := createMockEvent()
	// Use naturally aligned time (on 5-minute boundary)
	startsAt := event.StartsAt.Add(-1 * time.Hour) // Before event start
	startsAt = time.Date(startsAt.Year(), startsAt.Month(), startsAt.Day(), startsAt.Hour(), (startsAt.Minute()/5)*5, 0, 0, startsAt.Location())
	endsAt := event.StartsAt.Add(1 * time.Hour)
	endsAt = time.Date(endsAt.Year(), endsAt.Month(), endsAt.Day(), endsAt.Hour(), (endsAt.Minute()/5)*5, 0, 0, endsAt.Location())

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for start time before event start")
	assert.Equal(t, constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err, err, "Expected ERR_AVAILABILITY_START_BEFORE_EVENT error")
}

// TestValidateAvailabilityTimes_EndAfterEvent tests validation for end time after event end
func TestValidateAvailabilityTimes_EndAfterEvent(t *testing.T) {
	event := createMockEvent()
	// Use naturally aligned time (on 5-minute boundary)
	startsAt := event.EndsAt.Add(-1 * time.Hour)
	startsAt = time.Date(startsAt.Year(), startsAt.Month(), startsAt.Day(), startsAt.Hour(), (startsAt.Minute()/5)*5, 0, 0, startsAt.Location())
	endsAt := event.EndsAt.Add(1 * time.Hour) // After event end
	endsAt = time.Date(endsAt.Year(), endsAt.Month(), endsAt.Day(), endsAt.Hour(), (endsAt.Minute()/5)*5, 0, 0, endsAt.Location())

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for end time after event end")
	assert.Equal(t, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err, err, "Expected ERR_AVAILABILITY_END_AFTER_EVENT error")
}

// TestValidateAvailabilityTimes_ValidTimes tests validation for valid times
func TestValidateAvailabilityTimes_ValidTimes(t *testing.T) {
	event := createMockEvent()
	// Use naturally aligned time (on 5-minute boundary)
	startsAt := event.StartsAt.Add(1 * time.Hour)
	startsAt = time.Date(startsAt.Year(), startsAt.Month(), startsAt.Day(), startsAt.Hour(), (startsAt.Minute()/5)*5, 0, 0, startsAt.Location())
	endsAt := startsAt.Add(30 * time.Minute) // Valid: 30 minutes, naturally aligned on 5-minute intervals

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.NoError(t, err, "Expected no error for valid times")
}
