package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountEvent struct {
	AccountId uuid.UUID `gorm:"column:account_id;type:uuid;primaryKey" json:"-"`
	EventId   uuid.UUID `gorm:"column:event_id;type:uuid;primaryKey" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	// Relations
	Account Account `gorm:"foreignKey:AccountId;references:Id" json:"account"`
	Event   Event   `gorm:"foreignKey:EventId;references:Id" json:"event"`
}

func (AccountEvent) TableName() string {
	return "account_event"
}
