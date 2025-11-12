package event

import (
	model "app/db/models"
	"time"
)

type EventCreateDto struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	Duration    int       `json:"duration" binding:"required,min=15,max=30240"`
	StartsAt    time.Time `json:"startsAt" binding:"required"`
	EndsAt      time.Time `json:"endsAt" binding:"required"`
}

type EventResponse struct {
	model.Event
	Accounts []model.Account `json:"participants"`
}
