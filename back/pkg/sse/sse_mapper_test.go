package sse

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapToSSESlotUpdateMessages(t *testing.T) {
	startsAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	slots := []model.Slot{
		{Id: uuid.New(), EventId: uuid.New(), StartsAt: startsAt, EndsAt: startsAt.Add(time.Hour), IsValidated: true},
		{Id: uuid.New(), EventId: uuid.New(), StartsAt: startsAt.Add(2 * time.Hour), EndsAt: startsAt.Add(3 * time.Hour)},
	}

	messages := MapToSSESlotUpdateMessages(constants.SLOT_ACTION_UPDATED, slots)

	assert.Equal(t, []SSESlotUpdateMessage{
		{Action: constants.SLOT_ACTION_UPDATED, Id: slots[0].Id, StartsAt: slots[0].StartsAt, EndsAt: slots[0].EndsAt, IsValidated: true},
		{Action: constants.SLOT_ACTION_UPDATED, Id: slots[1].Id, StartsAt: slots[1].StartsAt, EndsAt: slots[1].EndsAt, IsValidated: false},
	}, messages)
}
