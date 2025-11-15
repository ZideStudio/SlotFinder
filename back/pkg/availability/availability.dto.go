package availability

import (
	"time"

	"github.com/google/uuid"
)

type AvailabilityCreateDto struct {
	EventId  uuid.UUID `json:"eventId" binding:"required,uuid"`
	StartsAt time.Time `json:"startsAt" binding:"required"`
	EndsAt   time.Time `json:"endsAt" binding:"required"`
}
