package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id          uuid.UUID `gorm:"column:id;type:uuid;unique;primary_key" json:"id"`
	Name        string    `gorm:"column:name;size:255" json:"name"`
	Description *string   `gorm:"column:description;type:text" json:"description"`
	Duration    int       `gorm:"column:duration;default:60" json:"duration"` // In minutes
	StartsAt    time.Time `gorm:"column:starts_at" json:"startsAt"`
	EndsAt      time.Time `gorm:"column:ends_at" json:"endsAt"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	OwnerId     uuid.UUID `gorm:"column:owner_id;type:uuid;primaryKey" json:"-"`

	// Relations
	Owner          Account        `gorm:"foreignKey:OwnerId;references:Id" json:"owner"`
	AccountEvents  []AccountEvent `gorm:"foreignKey:EventId;references:Id" json:"-"`
	Availabilities []Availability `gorm:"foreignKey:EventId;references:Id" json:"availabilities"`
}

func (Event) TableName() string {
	return "event"
}

func (e *Event) Sanitized() *Event {
	e.Owner = e.Owner.Sanitized()

	if len(e.Availabilities) != 0 {
		availabilities := make([]Availability, len(e.Availabilities))
		for i, availability := range e.Availabilities {
			availabilities[i] = *availability.Sanitized()
		}

		e.Availabilities = availabilities
	}

	return e
}

func (e *Event) HasUserAccess(userId *uuid.UUID) bool {
	if userId == nil {
		return false
	}

	var hasAccess bool
	for _, accountEvent := range e.AccountEvents {
		if accountEvent.AccountId == *userId {
			hasAccess = true
		}
	}

	return hasAccess
}

func (e *Event) HasEnded() bool {
	return time.Now().After(e.EndsAt)
}
