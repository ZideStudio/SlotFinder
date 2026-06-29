package repository

import (
	"app/db"
	model "app/db/models"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AccountProvidersRepository struct {
	db *gorm.DB
}

func NewAccountProvidersRepository(database *gorm.DB) *AccountProvidersRepository {
	if database == nil {
		database = db.GetDB()
	}
	return &AccountProvidersRepository{
		db: database,
	}
}

type AccountProviderCreateDto struct {
	UserName  string
	Email     string
	Password  string
	Providers []model.AccountProvider
}

func (r *AccountProvidersRepository) Create(accountProvider model.AccountProvider) error {
	if err := r.db.Create(&accountProvider).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::CREATE Failed to create account provider")
		return err
	}

	return nil
}

func (r *AccountProvidersRepository) FindOneById(id string, provider string, accountProvider *model.AccountProvider) error {
	err := r.db.Where("LOWER(provider) = LOWER(?) AND LOWER(id) = LOWER(?)", provider, id).Preload("Account").First(&accountProvider).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::FIND_ONE_BY_ID Failed to find account provider")
	}
	return err
}

func (r *AccountProvidersRepository) Delete(id string) error {
	err := r.db.Where("LOWER(id) = LOWER(?)", id).Delete(&model.AccountProvider{}).Error
	if err != nil {
		log.Error().Err(err).Msg("ACCOUNT_PROVIDERS_REPOSITORY::DELETE Failed to delete account provider")
	}
	return err
}
