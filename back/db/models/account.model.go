package model

import (
	"app/commons/constants"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id                   uuid.UUID                 `gorm:"column:id;type:uuid;unique;primary_key" json:"id,omitzero"`
	UserName             *string                   `gorm:"column:username;default:null;size:255" json:"userName"`
	Email                *string                   `gorm:"column:email;default:null;size:255" json:"email,omitempty"`
	Password             *string                   `gorm:"column:password;size:255" json:"-"`
	AvatarUrl            string                    `gorm:"column:avatar_url;size:255;default:null" json:"avatarUrl"`
	Language             constants.AccountLanguage `gorm:"column:language;type:VARCHAR(10);default:'en'" json:"language"`
	Color                string                    `gorm:"column:color;size:7" json:"color"`
	ResetToken           *string                   `gorm:"column:reset_token;size:255;default:null" json:"-"`
	PasswordResetTokenAt *time.Time                `gorm:"column:password_reset_token_at;default:null" json:"-"`
	Events               []AccountEvent            `gorm:"foreignKey:AccountId;references:Id" json:"events,omitempty"`
	Providers            []AccountProvider         `gorm:"foreignKey:AccountId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"providers,omitempty"`
	CreatedAt            time.Time                 `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt,omitzero"`
	UpdatedAt            time.Time                 `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt            *time.Time                `gorm:"column:deleted_at;default:null" json:"-"`
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

func (a *Account) Sanitized(overrideColor *string) Account {
	color := a.Color
	if overrideColor != nil {
		color = *overrideColor
	}

	return Account{
		UserName:  a.UserName,
		AvatarUrl: a.AvatarUrl,
		Color:     color,
	}
}
