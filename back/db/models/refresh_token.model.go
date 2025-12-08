package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id        uuid.UUID  `gorm:"column:id;type:uuid;unique;primary_key" json:"id,omitzero"`
	AccountId uuid.UUID  `gorm:"column:account_id;type:uuid;not null" json:"accountId"`
	TokenHash string     `gorm:"column:token_hash;type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	IsRevoked bool       `gorm:"column:is_revoked;default:false" json:"isRevoked"`
	CreatedAt time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt,omitzero"`
	UpdatedAt time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"-"`
	RevokedAt *time.Time `gorm:"column:revoked_at;default:null" json:"-"`
	Account   *Account   `gorm:"foreignKey:AccountId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"account,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
