package model

import (
	"app/commons/constants"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvent_Sanitized_WithAvailabilities(t *testing.T) {
	username := "alice"
	event := Event{
		Owner: Account{UserName: &username},
		Availabilities: []Availability{
			{Account: Account{UserName: &username}},
		},
	}

	sanitized := event.Sanitized()

	assert.Len(t, sanitized.Availabilities, 1)
	assert.Equal(t, "alice", sanitized.Availabilities[0].UserName)
}

func TestEvent_HasUserAccess_NilUserId(t *testing.T) {
	event := Event{}
	assert.False(t, event.HasUserAccess(nil))
}

func TestEvent_IsOwner_NilUserId(t *testing.T) {
	event := Event{}
	assert.False(t, event.IsOwner(nil))
}

func TestEvent_GetValidatedSlot_NoneValidated(t *testing.T) {
	event := Event{Slots: []Slot{{IsValidated: false}}}
	assert.Nil(t, event.GetValidatedSlot())
}

func TestEvent_CheckAndAutoUpdateStatus_ExpiredEvent_UpdateSucceeds(t *testing.T) {
	event := Event{
		Status: constants.EVENT_STATUS_UPCOMING,
		EndsAt: time.Now().Add(-time.Hour),
	}
	requireOneOfStatuses := []constants.EventStatus{constants.EVENT_STATUS_FINISHED}

	hasStatus, err := event.CheckAndAutoUpdateStatus(func(e *Event) error {
		return nil
	}, &requireOneOfStatuses)

	assert.NoError(t, err)
	assert.True(t, hasStatus)
	assert.Equal(t, constants.EVENT_STATUS_FINISHED, event.Status)
}

func TestEvent_CheckAndAutoUpdateStatus_ExpiredEvent_UpdateFails(t *testing.T) {
	event := Event{
		Status: constants.EVENT_STATUS_UPCOMING,
		EndsAt: time.Now().Add(-time.Hour),
	}

	hasStatus, err := event.CheckAndAutoUpdateStatus(func(e *Event) error {
		return errors.New("db unavailable")
	}, nil)

	assert.Error(t, err)
	assert.False(t, hasStatus)
}

func TestEvent_CheckAndAutoUpdateStatus_NotExpired_SkipsUpdate(t *testing.T) {
	event := Event{
		Status: constants.EVENT_STATUS_UPCOMING,
		EndsAt: time.Now().Add(time.Hour),
	}
	requireOneOfStatuses := []constants.EventStatus{constants.EVENT_STATUS_UPCOMING}

	called := false
	hasStatus, err := event.CheckAndAutoUpdateStatus(func(e *Event) error {
		called = true
		return nil
	}, &requireOneOfStatuses)

	assert.NoError(t, err)
	assert.True(t, hasStatus)
	assert.False(t, called)
	assert.Equal(t, constants.EVENT_STATUS_UPCOMING, event.Status)
}

func TestEvent_HasUserAccess_Found(t *testing.T) {
	userId := uuid.New()
	event := Event{AccountEvents: []AccountEvent{{AccountId: userId}}}
	assert.True(t, event.HasUserAccess(&userId))
}
