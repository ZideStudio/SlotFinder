package model

import (
	"time"

	"github.com/google/uuid"
)

type Availability struct {
	Id        uuid.UUID `gorm:"column:id;type:uuid;unique;primary_key" json:"id,omitzero"`
	AccountId uuid.UUID `gorm:"column:account_id;type:uuid;primaryKey" json:"-"`
	Username  string    `gorm:"-" json:"username"`
	Account   Account   `gorm:"foreignKey:AccountId;references:Id" json:"-"`
	EventId   uuid.UUID `gorm:"column:event_id;type:uuid;primaryKey" json:"-"`
	Event     Event     `gorm:"foreignKey:EventId;references:Id" json:"-"`
	StartsAt  time.Time `gorm:"column:starts_at" json:"startsAt"`
	EndsAt    time.Time `gorm:"column:ends_at" json:"endsAt"`
}

func (Availability) TableName() string {
	return "availability"
}

func (a *Availability) Sanitized() *Availability {
	if a.Account.Username != nil {
		a.Username = *a.Account.Username
	}
	return a
}
