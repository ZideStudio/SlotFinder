package event

import (
	"app/commons/constants"
	model "app/db/models"
	"time"

	"github.com/google/uuid"
)

// EventDurationFields - duration split into days/hours/minutes
type EventDurationFields struct {
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// EventOwnerDto - owner with event-specific color
type EventOwnerDto struct {
	UserName  *string `json:"userName"`
	AvatarUrl string  `json:"avatarUrl"`
	Color     string  `json:"color"`
}

// EventParticipantDto - participant with event-specific color
type EventParticipantDto struct {
	UserName  *string `json:"userName"`
	AvatarUrl string  `json:"avatarUrl"`
	Color     string  `json:"color"`
}

// EventListItemDto - GET /events (paginated, no joins)
type EventListItemDto struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	EventDurationFields
	StartsAt time.Time            `json:"startsAt"`
	EndsAt   time.Time            `json:"endsAt"`
	Status   constants.EventStatus `json:"status"`
}

// EventCreateResponseDto - POST /events (event + owner)
type EventCreateResponseDto struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	EventDurationFields
	StartsAt time.Time            `json:"startsAt"`
	EndsAt   time.Time            `json:"endsAt"`
	Status   constants.EventStatus `json:"status"`
	Owner    EventOwnerDto        `json:"owner"`
}

// EventBasicResponseDto - GET /events/:id/summary (public)
type EventBasicResponseDto struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	EventDurationFields
	StartsAt time.Time            `json:"startsAt"`
	EndsAt   time.Time            `json:"endsAt"`
	Status   constants.EventStatus `json:"status"`
}

// EventFullResponseDto - GET /events/:id (member) and POST /events/:id/join
type EventFullResponseDto struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	EventDurationFields
	StartsAt       time.Time             `json:"startsAt"`
	EndsAt         time.Time             `json:"endsAt"`
	Status         constants.EventStatus `json:"status"`
	Owner          EventOwnerDto         `json:"owner"`
	Participants   []EventParticipantDto `json:"participants"`
	Availabilities []model.Availability  `json:"availabilities"`
	Slots          []model.Slot          `json:"slots"`
}
