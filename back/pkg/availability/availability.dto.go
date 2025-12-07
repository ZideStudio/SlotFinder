package availability

import (
	"time"
)

type AvailabilityCreateDto struct {
	StartsAt time.Time `json:"startsAt" binding:"required"`
	EndsAt   time.Time `json:"endsAt" binding:"required"`
}

type AvailabilityUpdateDto struct {
	StartsAt time.Time `json:"startsAt" binding:"required"`
	EndsAt   time.Time `json:"endsAt" binding:"required"`
}
