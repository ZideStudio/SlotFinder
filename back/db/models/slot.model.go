package model

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	Id          uuid.UUID `gorm:"column:id;type:uuid;unique;primary_key" json:"id,omitzero"`
	EventId     uuid.UUID `gorm:"column:event_id;type:uuid;primaryKey" json:"-"`
	Event       Event     `gorm:"foreignKey:EventId;references:Id" json:"-"`
	StartsAt    time.Time `gorm:"column:starts_at" json:"startsAt"`
	EndsAt      time.Time `gorm:"column:ends_at" json:"endsAt"`
	IsValidated bool      `gorm:"column:is_validated;default:false" json:"isValidated"`
}

func (Slot) TableName() string {
	return "slot"
}
