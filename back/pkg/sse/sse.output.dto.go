package sse

import (
	"time"

	"github.com/google/uuid"
)

// SSESlotUpdateMessage is the payload sent on each SSE data frame: an array of computed slots
type SSESlotUpdateMessage struct {
	Id          uuid.UUID `json:"id"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	IsValidated bool      `json:"isValidated"`
}
