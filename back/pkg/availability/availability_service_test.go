package availability

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// NOTE: Do not initialize real services at package init time in unit tests.
// NewAvailabilityService(nil) wires real dependencies indirectly (via other services),
// which can crash tests when configuration isn't set.
// Create the service inside each test instead.

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

// Helper function to align time to 5-minute boundaries
func alignToFiveMinutes(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), (t.Minute()/5)*5, 0, 0, t.Location())
}

// TestValidateAvailabilityTimes_StartAfterEnd tests validation for end date before start date
func TestValidateAvailabilityTimes_StartAfterEnd(t *testing.T) {
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

	event := createMockEvent()
	startsAt := event.StartsAt.Add(2 * time.Hour)
	endsAt := event.StartsAt.Add(1 * time.Hour) // End before start

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for end date before start date")
	assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err, "Expected ERR_EVENT_START_AFTER_END error")
}

// TestValidateAvailabilityTimes_DurationTooShort tests validation for duration less than 5 minutes
func TestValidateAvailabilityTimes_DurationTooShort(t *testing.T) {
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

	event := createMockEvent()
	startsAt := event.StartsAt.Add(1 * time.Hour)
	endsAt := startsAt.Add(3 * time.Minute) // Less than 5 minutes

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for duration less than 5 minutes")
	assert.Equal(t, constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err, err, "Expected ERR_AVAILABILITY_DURATION_TOO_SHORT error")
}

// TestValidateAvailabilityTimes_InvalidTimeInterval tests validation for times not aligned on 5-minute intervals
func TestValidateAvailabilityTimes_InvalidTimeInterval(t *testing.T) {
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

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
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

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
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

	event := createMockEvent()
	// Use helper to align time to 5-minute boundary
	startsAt := alignToFiveMinutes(event.StartsAt.Add(-1 * time.Hour)) // Before event start
	endsAt := alignToFiveMinutes(event.StartsAt.Add(1 * time.Hour))

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for start time before event start")
	assert.Equal(t, constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err, err, "Expected ERR_AVAILABILITY_START_BEFORE_EVENT error")
}

// TestValidateAvailabilityTimes_EndAfterEvent tests validation for end time after event end
func TestValidateAvailabilityTimes_EndAfterEvent(t *testing.T) {
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

	event := createMockEvent()
	// Use helper to align time to 5-minute boundary
	startsAt := alignToFiveMinutes(event.EndsAt.Add(-1 * time.Hour))
	endsAt := alignToFiveMinutes(event.EndsAt.Add(1 * time.Hour)) // After event end

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.Error(t, err, "Expected error for end time after event end")
	assert.Equal(t, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err, err, "Expected ERR_AVAILABILITY_END_AFTER_EVENT error")
}

// TestValidateAvailabilityTimes_ValidTimes tests validation for valid times
func TestValidateAvailabilityTimes_ValidTimes(t *testing.T) {
	// Minimal service: validateAvailabilityTimes doesn't require external dependencies.
	service := &AvailabilityService{}

	event := createMockEvent()
	// Use helper to align time to 5-minute boundary
	startsAt := alignToFiveMinutes(event.StartsAt.Add(1 * time.Hour))
	endsAt := startsAt.Add(30 * time.Minute) // Valid: 30 minutes, naturally aligned on 5-minute intervals

	err := service.validateAvailabilityTimes(startsAt, endsAt, &event)
	assert.NoError(t, err, "Expected no error for valid times")
}
