package event

import (
	model "app/db/models"
	"time"
)

type EventCreateDto struct {
	Name     string    `json:"name" binding:"required"`
	Duration int       `json:"duration" binding:"required,min=15,max=30240"`
	StartsAt time.Time `json:"starts_at" binding:"required"`
	EndsAt   time.Time `json:"ends_at" binding:"required"`
}

type EventResponse struct {
	model.Event
	Accounts []model.Account `json:"accounts"`
}
