package availability

import (
	"time"

	"github.com/google/uuid"
)

// AvailabilityResponseDto - POST /events/:id/availability and PATCH /availabilities/:id
type AvailabilityResponseDto struct {
	Id       uuid.UUID `json:"id"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}
