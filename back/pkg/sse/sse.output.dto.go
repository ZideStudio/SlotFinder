package sse

import (
	"app/commons/constants"
	"time"

	"github.com/google/uuid"
)

// SSESlotUpdateMessage represents one slot change sent in SSE data frames (each "data" payload is a JSON array of these,
// ordered DELETED, UPDATED then CREATED). No snapshot is sent on connection: current slots come from GET /events/{eventId}.
type SSESlotUpdateMessage struct {
	Action      constants.SlotAction `json:"action"`
	Id          uuid.UUID            `json:"id"`
	StartsAt    time.Time            `json:"startsAt"`
	EndsAt      time.Time            `json:"endsAt"`
	IsValidated bool                 `json:"isValidated"`
}
