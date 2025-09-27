package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/rs/zerolog/log"
)

type AccountEventRepository struct{}

func (*AccountEventRepository) Create(event *model.AccountEvent) error {
	if err := db.GetDB().Preload("Account").Create(&event).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::CREATE Failed to create account")
		return err
	}

	return nil
}
