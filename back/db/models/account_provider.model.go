package model

import (
	"app/commons/constants"

	"github.com/google/uuid"
)

type AccountProvider struct {
	AccountId uuid.UUID          `json:"-" gorm:"primaryKey;type:uuid;column:-"`
	Provider  constants.Provider `json:"provider" gorm:"primaryKey;type:varchar(20);column:provider"`
	Id        string             `json:"-" gorm:"type:varchar(100);column:-"`
	Email     string             `json:"-" gorm:"type:varchar(100);column:-"`
	Username  string             `json:"-" gorm:"type:varchar(100);column:-"`
}

func (AccountProvider) TableName() string {
	return "account_provider"
}
