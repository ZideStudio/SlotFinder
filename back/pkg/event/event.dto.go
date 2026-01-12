package event

import (
	"app/commons/constants"
	"time"
)

type EventCreateDto struct {
	Name        string    `json:"name" binding:"required,min=5,max=100"`
	Description *string   `json:"description" binding:"omitempty,max=500"`
	Duration    int       `json:"duration" binding:"required,min=15,max=30240"`
	StartsAt    time.Time `json:"startsAt" binding:"required"`
	EndsAt      time.Time `json:"endsAt" binding:"required"`
}

type EventUpdateDto struct {
	Name        *string                `json:"name" binding:"omitempty,min=5,max=100"`
	Description *string                `json:"description" binding:"omitempty,max=500"`
	Duration    *int                   `json:"duration" binding:"omitempty,min=15,max=30240"`
	StartsAt    *time.Time             `json:"startsAt"`
	EndsAt      *time.Time             `json:"endsAt"`
	Status      *constants.EventStatus `json:"status" binding:"omitempty,oneof=IN_DECISION"`
}

type EventProfileDto struct {
	Color string `json:"color"`
}
