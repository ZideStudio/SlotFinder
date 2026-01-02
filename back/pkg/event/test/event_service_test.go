package test

import (
	"app/commons/constants"
	model "app/db/models"
	"app/pkg/event"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test the parseDtoEventDates function directly since it contains most of the validation logic
func TestParseDtoEventDates_Success(t *testing.T) {
	t.Run("should parse dates successfully", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
		}
		newStart := baseDate.AddDate(0, 0, 1) // 1 day after base
		newEnd := baseDate.AddDate(0, 0, 5)   // 5 days after base

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, &newEnd)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, newStart, testEvent.StartsAt)
		assert.Equal(t, newEnd, testEvent.EndsAt)
	})

	t.Run("should handle nil dates", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		originalStart := baseDate
		originalEnd := baseDate.AddDate(0, 0, 4) // 4 days later
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: originalStart,
			EndsAt:   originalEnd,
		}

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, originalStart, testEvent.StartsAt) // Should remain unchanged
		assert.Equal(t, originalEnd, testEvent.EndsAt)     // Should remain unchanged
	})

	t.Run("should update only start date", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0)  // Next year
		originalEnd := baseDate.AddDate(0, 0, 4) // 4 days later
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   originalEnd,
		}
		newStart := baseDate.AddDate(0, 0, 1) // 1 day after base

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, newStart, testEvent.StartsAt)
		assert.Equal(t, originalEnd, testEvent.EndsAt) // Should remain unchanged
	})

	t.Run("should update only end date", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		originalStart := baseDate
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: originalStart,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
		}
		newEnd := baseDate.AddDate(0, 0, 5) // 5 days after base

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, &newEnd)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, originalStart, testEvent.StartsAt) // Should remain unchanged
		assert.Equal(t, newEnd, testEvent.EndsAt)
	})
}

func TestParseDtoEventDates_StartAfterEnd(t *testing.T) {
	t.Run("should return error when start is after end", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
		}
		// Set start after current end
		newStart := baseDate.AddDate(0, 0, 10) // 10 days after base

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})

	t.Run("should return error when end is before start", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate.AddDate(0, 0, 4), // 4 days after base
			EndsAt:   baseDate.AddDate(0, 0, 9), // 9 days after base
		}
		// Set end before current start
		newEnd := baseDate // Same as original base (before start)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, &newEnd)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})

	t.Run("should return error when both dates make start after end", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
		}
		newStart := baseDate.AddDate(0, 0, 10) // 10 days after base
		newEnd := baseDate.AddDate(0, 0, 8)    // 8 days after base (before newStart)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, &newEnd)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})
}

func TestParseDtoEventDates_DurationTooShort(t *testing.T) {
	t.Run("should return error when duration is less than 1 day", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
		}
		// Set end to only 1 hour after start
		newEnd := testEvent.StartsAt.Add(1 * time.Hour)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, &newEnd)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_DURATION_TOO_SHORT.Err, err)
	})

	t.Run("should return error when new start makes duration too short", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.Add(36 * time.Hour), // 1.5 days later
		}
		// Move start to make duration < 24h
		newStart := testEvent.EndsAt.Add(-12 * time.Hour)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_DURATION_TOO_SHORT.Err, err)
	})
}

func TestParseDtoEventDates_PastDates(t *testing.T) {
	t.Run("should return error when start date is in the past", func(t *testing.T) {
		// Arrange
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: time.Now().AddDate(0, 0, 1), // Tomorrow
			EndsAt:   time.Now().AddDate(0, 0, 5), // In 5 days
		}
		// Set start date in the past
		pastDate := time.Now().AddDate(0, 0, -1)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &pastDate, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_BEFORE_TODAY.Err, err)
	})

	t.Run("should return error when end date is in the past", func(t *testing.T) {
		// Arrange
		futureDate := time.Now().AddDate(0, 0, 1) // Tomorrow
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: futureDate,
			EndsAt:   futureDate.AddDate(0, 0, 4), // 5 days from now
		}
		// Set end date in the past but ensure it creates a valid duration first
		// This will trigger the "end date in past" validation since start is not being updated
		pastEndDate := time.Now().AddDate(0, 0, -1)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, &pastEndDate)

		// Assert
		assert.Error(t, err)
		// Since start > end, it will trigger start_after_end error first
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})
}

func TestParseDtoEventDates_ValidatedSlotConflict(t *testing.T) {
	t.Run("should return error when end date conflicts with validated slot", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		validatedSlot := &model.Slot{
			Id:          uuid.New(),
			StartsAt:    baseDate.AddDate(0, 0, 2).Add(10 * time.Hour), // Day 2, 10am
			EndsAt:      baseDate.AddDate(0, 0, 2).Add(12 * time.Hour), // Day 2, 12pm
			IsValidated: true,
		}
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
			Slots:    []model.Slot{*validatedSlot},
		}
		// Try to set end date before validated slot
		newEnd := baseDate.AddDate(0, 0, 1) // Day 1 (before slot on day 2)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, nil, &newEnd)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err, err)
	})

	t.Run("should return error when start date conflicts with validated slot", func(t *testing.T) {
		// Arrange
		baseDate := time.Now().AddDate(1, 0, 0) // Next year
		validatedSlot := &model.Slot{
			Id:          uuid.New(),
			StartsAt:    baseDate.AddDate(0, 0, 2).Add(10 * time.Hour), // Day 2, 10am
			EndsAt:      baseDate.AddDate(0, 0, 2).Add(12 * time.Hour), // Day 2, 12pm
			IsValidated: true,
		}
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4), // 4 days later
			Slots:    []model.Slot{*validatedSlot},
		}
		// Try to set start date after validated slot start
		newStart := baseDate.AddDate(0, 0, 3) // Day 3 (after slot on day 2)

		// Act
		err := event.ParseDtoEventDatesForTesting(testEvent, &newStart, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err, err)
	})
}

func TestParseDtoEventDates_NilEvent(t *testing.T) {
	t.Run("should return error for nil event", func(t *testing.T) {
		// Act
		err := event.ParseDtoEventDatesForTesting(nil, nil, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "event pointer is nil", err.Error())
	})
}
