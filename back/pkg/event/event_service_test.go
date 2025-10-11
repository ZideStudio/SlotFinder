package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var service = NewEventService(nil)

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
		Duration: 60,
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
		Duration: 60,
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
		Duration: 60,
		StartsAt: yesterday,
		EndsAt:   twoDaysLater,
	}

	_, err := service.Create(data, user)
	assert.Error(t, err, "Expected error for event starting in the past")
	assert.Equal(t, constants.ERR_EVENT_START_BEFORE_TODAY.Err, err, "Expected ERR_EVENT_START_BEFORE_TODAY error")
}
