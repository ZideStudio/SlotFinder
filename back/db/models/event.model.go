package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id        uuid.UUID `gorm:"column:id;type:uuid;unique;primary_key" json:"id"`
	Name      string    `gorm:"column:name;size:255" json:"name"`
	Duration  int       `gorm:"column:duration;default:60" json:"duration"` // In minutes
	StartsAt  time.Time `gorm:"column:starts_at" json:"startsAt"`
	EndsAt    time.Time `gorm:"column:ends_at" json:"endsAt"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	OwnerId   uuid.UUID `gorm:"column:owner_id;type:uuid;primaryKey" json:"-"`

	// Relations
	Owner         Account        `gorm:"foreignKey:OwnerId;references:Id" json:"owner"`
	AccountEvents []AccountEvent `gorm:"foreignKey:EventId;references:Id" json:"-"`
}

func (Event) TableName() string {
	return "event"
}
