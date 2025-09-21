package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/rs/zerolog/log"
)

type AccountProvidersRepository struct{}

type AccountProviderCreateDto struct {
	UserName  string
	Email     string
	Password  string
	Providers []model.AccountProvider
}

func (*AccountProvidersRepository) Create(accountProvider model.AccountProvider) error {
	if err := db.GetDB().Create(&accountProvider).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::CREATE Failed to create account provider")
		return err
	}

	return nil
}

func (*AccountProvidersRepository) FindOneById(id string, provider string, accountProvider *model.AccountProvider) error {
	err := db.GetDB().Where("LOWER(provider) = LOWER(?) AND LOWER(id) = LOWER(?)", provider, id).First(&accountProvider).Error
	if err != nil {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::FIND_ONE_BY_ID Failed to find account provider")
	}
	return err
}

func (*AccountProvidersRepository) Delete(id string) error {
	err := db.GetDB().Where("LOWER(id) = LOWER(?)", id).Delete(&model.AccountProvider{}).Error
	if err != nil {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::DELETE Failed to delete account provider")
	}
	return err
}
