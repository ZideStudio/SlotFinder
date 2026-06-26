package slot

import (
	"time"

	"github.com/google/uuid"
)

// SlotResponseDto - POST /slots/:id/confirm
type SlotResponseDto struct {
	Id          uuid.UUID `json:"id"`
	IsValidated bool      `json:"isValidated"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
}
