package event

import "time"

type EventCreateDto struct {
	Name     string    `json:"name" binding:"required"`
	StartsAt time.Time `json:"starts_at" binding:"required"`
	EndsAt   time.Time `json:"ends_at" binding:"required"`
}
