package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id                   uuid.UUID         `gorm:"column:id;type:uuid;unique;primary_key" json:"id"`
	UserName             string            `gorm:"column:username;size:255" json:"userName"`
	Email                string            `gorm:"column:email;default:null;size:255" json:"email"`
	Password             *string           `gorm:"column:password;size:255" json:"-"`
	ResetToken           *string           `gorm:"column:reset_token;size:255;default:null" json:"-"`
	PasswordResetTokenAt *time.Time        `gorm:"column:password_reset_token_at;default:null" json:"-"`
	Events               []AccountEvent    `gorm:"foreignKey:AccountId;references:Id" json:"events"`
	Providers            []AccountProvider `gorm:"foreignKey:AccountId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"providers"`
	CreatedAt            time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt            time.Time         `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt            *time.Time        `gorm:"column:deleted_at;default:null" json:"-"`
}

func (Account) TableName() string {
	return "account"
}

func (a *Account) ComparePassword(password string) bool {
	if a.Password == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(*a.Password), []byte(password))
	return err == nil
}
