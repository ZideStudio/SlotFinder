package model

import (
	"app/commons/constants"

	"github.com/google/uuid"
)

type AccountProvider struct {
	AccountId uuid.UUID          `gorm:"primaryKey;type:uuid;column:account_id" json:"-"`
	Provider  constants.Provider `gorm:"primaryKey;type:varchar(20);column:provider" json:"provider"`
	Id        string             `gorm:"type:varchar(100);column:id" json:"-"`
}

func (AccountProvider) TableName() string {
	return "account_provider"
}
