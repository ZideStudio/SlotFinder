package event

import "time"

// EventCreateDto - POST /events
type EventCreateDto struct {
	Name        string    `json:"name" binding:"required,min=5,max=100"`
	Description *string   `json:"description" binding:"omitempty,max=500"`
	Days        int       `json:"days" binding:"min=0"`
	Hours       int       `json:"hours" binding:"min=0,max=23"`
	Minutes     int       `json:"minutes" binding:"min=0,max=59"`
	StartsAt    time.Time `json:"startsAt" binding:"required"`
	EndsAt      time.Time `json:"endsAt" binding:"required"`
}

// EventUpdateDto - PATCH /events/:id
type EventUpdateDto struct {
	Name        *string    `json:"name" binding:"omitempty,min=5,max=100"`
	Description *string    `json:"description" binding:"omitempty,max=500"`
	Days        *int       `json:"days" binding:"omitempty,min=0"`
	Hours       *int       `json:"hours" binding:"omitempty,min=0,max=23"`
	Minutes     *int       `json:"minutes" binding:"omitempty,min=0,max=59"`
	StartsAt    *time.Time `json:"startsAt"`
	EndsAt      *time.Time `json:"endsAt"`
}

// EventProfileDto - PATCH /events/:id/profile
type EventProfileDto struct {
	Color string `json:"color"`
}
