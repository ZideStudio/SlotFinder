package sse

import (
	"time"

	"github.com/google/uuid"
)

// SSESlotUpdateMessage represents one slot entry sent in SSE data frames (each "data" payload is a JSON array of these).
type SSESlotUpdateMessage struct {
	Id          uuid.UUID `json:"id"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	IsValidated bool      `json:"isValidated"`
}
