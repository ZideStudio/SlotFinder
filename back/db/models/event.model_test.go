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
	t.Parallel()
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
	t.Parallel()
	event := Event{}
	assert.False(t, event.HasUserAccess(nil))
}

func TestEvent_IsOwner_NilUserId(t *testing.T) {
	t.Parallel()
	event := Event{}
	assert.False(t, event.IsOwner(nil))
}

func TestEvent_GetValidatedSlot_NoneValidated(t *testing.T) {
	t.Parallel()
	event := Event{Slots: []Slot{{IsValidated: false}}}
	assert.Nil(t, event.GetValidatedSlot())
}

func TestEvent_CheckAndAutoUpdateStatus_ExpiredEvent_UpdateSucceeds(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	userId := uuid.New()
	event := Event{AccountEvents: []AccountEvent{{AccountId: userId}}}
	assert.True(t, event.HasUserAccess(&userId))
}

func TestEvent_HasUserAccess_NotFound(t *testing.T) {
	t.Parallel()
	userId := uuid.New()
	event := Event{AccountEvents: []AccountEvent{{AccountId: uuid.New()}}}
	assert.False(t, event.HasUserAccess(&userId))
}

func TestEvent_IsOwner_Matching(t *testing.T) {
	t.Parallel()
	ownerId := uuid.New()
	event := Event{OwnerId: ownerId}
	assert.True(t, event.IsOwner(&ownerId))
}

func TestEvent_IsOwner_NotMatching(t *testing.T) {
	t.Parallel()
	event := Event{OwnerId: uuid.New()}
	otherId := uuid.New()
	assert.False(t, event.IsOwner(&otherId))
}

func TestEvent_GetValidatedSlot_EmptySlots(t *testing.T) {
	t.Parallel()
	event := Event{Slots: []Slot{}}
	assert.Nil(t, event.GetValidatedSlot())
}

func TestEvent_GetValidatedSlot_Found(t *testing.T) {
	t.Parallel()
	slotId := uuid.New()
	event := Event{Slots: []Slot{{IsValidated: false}, {Id: slotId, IsValidated: true}}}

	slot := event.GetValidatedSlot()

	assert.NotNil(t, slot)
	assert.Equal(t, slotId, slot.Id)
}

func TestEvent_HasOneOfStatuses_NilStatuses(t *testing.T) {
	t.Parallel()
	event := Event{Status: constants.EVENT_STATUS_UPCOMING}
	assert.False(t, event.HasOneOfStatuses(nil))
}

func TestEvent_HasOneOfStatuses_Matching(t *testing.T) {
	t.Parallel()
	event := Event{Status: constants.EVENT_STATUS_UPCOMING}
	statuses := []constants.EventStatus{constants.EVENT_STATUS_UPCOMING}
	assert.True(t, event.HasOneOfStatuses(&statuses))
}

func TestEvent_HasOneOfStatuses_NotMatching(t *testing.T) {
	t.Parallel()
	event := Event{Status: constants.EVENT_STATUS_UPCOMING}
	statuses := []constants.EventStatus{constants.EVENT_STATUS_FINISHED}
	assert.False(t, event.HasOneOfStatuses(&statuses))
}

func TestEvent_TableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "event", Event{}.TableName())
}

func TestSlot_TableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "slot", Slot{}.TableName())
}
