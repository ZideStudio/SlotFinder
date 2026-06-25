package event

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// NOTE: These tests only validate `Create` input validation.
// Avoid constructing `EventService` here: `NewEventService(nil)` wires real deps (mail/templates/config)
// that can crash during test init when config isn't set.

// Keep service creation minimal for unit tests: provide a zero-value service since `Create` validation
// doesn't rely on repositories/services.
var service = &EventService{}

var username = "testuser"
var user = &guard.Claims{
	Id:       uuid.New(),
	Username: &username,
}

func TestCreate_EventDurationTooShort(t *testing.T) {
	now := time.Now().UTC()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	almostTwoDaysLater := tomorrow.Add(23*time.Hour + 59*time.Minute)

	data := &EventCreateDto{
		Name:     "Test Event",
		Days:     0,
		Hours:    1,
		Minutes:  0,
		StartsAt: tomorrow,
		EndsAt:   almostTwoDaysLater,
	}

	_, err := service.Create(data, user)
	assert.Error(t, err, "Expected error for event duration less than 1 day")
	assert.Equal(t, constants.ERR_EVENT_DURATION_TOO_SHORT.Err, err, "Expected ERR_EVENT_DURATION_TOO_SHORT error")
}

func TestCreate_EventStartAfterEnd(t *testing.T) {
	now := time.Now().UTC()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)

	data := &EventCreateDto{
		Name:     "Test Event",
		Days:     0,
		Hours:    1,
		Minutes:  0,
		StartsAt: tomorrow,
		EndsAt:   yesterday,
	}

	_, err := service.Create(data, user)
	assert.Error(t, err, "Expected error for event with end date before start date")
	assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err, "Expected ERR_EVENT_START_AFTER_END error")
}

func TestCreate_EventStartInPast(t *testing.T) {
	now := time.Now().UTC()
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)
	twoDaysLater := yesterday.Add(48 * time.Hour)

	data := &EventCreateDto{
		Name:     "Test Event",
		Days:     0,
		Hours:    1,
		Minutes:  0,
		StartsAt: yesterday,
		EndsAt:   twoDaysLater,
	}

	_, err := service.Create(data, user)
	assert.Error(t, err, "Expected error for event starting in the past")
	assert.Equal(t, constants.ERR_EVENT_START_BEFORE_TODAY.Err, err, "Expected ERR_EVENT_START_BEFORE_TODAY error")
}

func TestParseDtoEventDates_Success(t *testing.T) {
	t.Run("should parse dates successfully", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
		}
		newStart := baseDate.AddDate(0, 0, 1)
		newEnd := baseDate.AddDate(0, 0, 5)

		err := SetEventDatesFromDto(testEvent, &newStart, &newEnd)

		assert.NoError(t, err)
		assert.Equal(t, newStart, testEvent.StartsAt)
		assert.Equal(t, newEnd, testEvent.EndsAt)
	})

	t.Run("should handle nil dates", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		originalStart := baseDate
		originalEnd := baseDate.AddDate(0, 0, 4)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: originalStart,
			EndsAt:   originalEnd,
		}

		err := SetEventDatesFromDto(testEvent, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, originalStart, testEvent.StartsAt)
		assert.Equal(t, originalEnd, testEvent.EndsAt)
	})

	t.Run("should update only start date", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		originalEnd := baseDate.AddDate(0, 0, 4)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   originalEnd,
		}
		newStart := baseDate.AddDate(0, 0, 1)

		err := SetEventDatesFromDto(testEvent, &newStart, nil)

		assert.NoError(t, err)
		assert.Equal(t, newStart, testEvent.StartsAt)
		assert.Equal(t, originalEnd, testEvent.EndsAt)
	})

	t.Run("should update only end date", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		originalStart := baseDate
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: originalStart,
			EndsAt:   baseDate.AddDate(0, 0, 4),
		}
		newEnd := baseDate.AddDate(0, 0, 5)

		err := SetEventDatesFromDto(testEvent, nil, &newEnd)

		assert.NoError(t, err)
		assert.Equal(t, originalStart, testEvent.StartsAt)
		assert.Equal(t, newEnd, testEvent.EndsAt)
	})
}

func TestParseDtoEventDates_StartAfterEnd(t *testing.T) {
	t.Run("should return error when start is after end", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
		}
		newStart := baseDate.AddDate(0, 0, 10)

		err := SetEventDatesFromDto(testEvent, &newStart, nil)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})

	t.Run("should return error when end is before start", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate.AddDate(0, 0, 4),
			EndsAt:   baseDate.AddDate(0, 0, 9),
		}
		newEnd := baseDate

		err := SetEventDatesFromDto(testEvent, nil, &newEnd)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})

	t.Run("should return error when both dates make start after end", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
		}
		newStart := baseDate.AddDate(0, 0, 10)
		newEnd := baseDate.AddDate(0, 0, 8)

		err := SetEventDatesFromDto(testEvent, &newStart, &newEnd)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})
}

func TestParseDtoEventDates_DurationTooShort(t *testing.T) {
	t.Run("should return error when duration is less than 1 day", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
		}
		newEnd := testEvent.StartsAt.Add(1 * time.Hour)

		err := SetEventDatesFromDto(testEvent, nil, &newEnd)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_DURATION_TOO_SHORT.Err, err)
	})

	t.Run("should return error when new start makes duration too short", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.Add(36 * time.Hour),
		}
		newStart := testEvent.EndsAt.Add(-12 * time.Hour)

		err := SetEventDatesFromDto(testEvent, &newStart, nil)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_DURATION_TOO_SHORT.Err, err)
	})
}

func TestParseDtoEventDates_PastDates(t *testing.T) {
	t.Run("should return error when start date is in the past", func(t *testing.T) {
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: time.Now().AddDate(0, 0, 1),
			EndsAt:   time.Now().AddDate(0, 0, 5),
		}
		pastDate := time.Now().AddDate(0, 0, -1)

		err := SetEventDatesFromDto(testEvent, &pastDate, nil)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_BEFORE_TODAY.Err, err)
	})

	t.Run("should return error when end date is in the past", func(t *testing.T) {
		futureDate := time.Now().AddDate(0, 0, 1)
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: futureDate,
			EndsAt:   futureDate.AddDate(0, 0, 4),
		}
		pastEndDate := time.Now().AddDate(0, 0, -1)

		err := SetEventDatesFromDto(testEvent, nil, &pastEndDate)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_EVENT_START_AFTER_END.Err, err)
	})
}

func TestParseDtoEventDates_ValidatedSlotConflict(t *testing.T) {
	t.Run("should return error when end date conflicts with validated slot", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		validatedSlot := &model.Slot{
			Id:          uuid.New(),
			StartsAt:    baseDate.AddDate(0, 0, 2).Add(10 * time.Hour),
			EndsAt:      baseDate.AddDate(0, 0, 2).Add(12 * time.Hour),
			IsValidated: true,
		}
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
			Slots:    []model.Slot{*validatedSlot},
		}
		newEnd := baseDate.AddDate(0, 0, 1).Add(1 * time.Hour)

		err := SetEventDatesFromDto(testEvent, nil, &newEnd)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err, err)
	})

	t.Run("should return error when start date conflicts with validated slot", func(t *testing.T) {
		baseDate := time.Now().AddDate(1, 0, 0)
		validatedSlot := &model.Slot{
			Id:          uuid.New(),
			StartsAt:    baseDate.AddDate(0, 0, 2).Add(10 * time.Hour),
			EndsAt:      baseDate.AddDate(0, 0, 2).Add(12 * time.Hour),
			IsValidated: true,
		}
		testEvent := &model.Event{
			Id:       uuid.New(),
			StartsAt: baseDate,
			EndsAt:   baseDate.AddDate(0, 0, 4),
			Slots:    []model.Slot{*validatedSlot},
		}
		newStart := baseDate.AddDate(0, 0, 3)

		err := SetEventDatesFromDto(testEvent, &newStart, nil)

		assert.Error(t, err)
		assert.Equal(t, constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err, err)
	})
}

func TestParseDtoEventDates_NilEvent(t *testing.T) {
	t.Run("should return error for nil event", func(t *testing.T) {
		err := SetEventDatesFromDto(nil, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, "event pointer is nil", err.Error())
	})
}
