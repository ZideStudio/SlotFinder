package repository

import (
	"app/db"
	model "app/db/models"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountEventRepository struct {
	db *gorm.DB
}

func NewAccountEventRepository(database *gorm.DB) *AccountEventRepository {
	if database == nil {
		database = db.GetDB()
	}
	return &AccountEventRepository{
		db: database,
	}
}

func (r *AccountEventRepository) Create(accountEvent *model.AccountEvent) error {
	if err := r.db.Preload("Account").Create(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::CREATE Failed to create account_event")
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (r *AccountEventRepository) Updates(accountEvent *model.AccountEvent) error {
	if err := r.db.Where("account_id = ? AND event_id = ?", accountEvent.AccountId, accountEvent.EventId).Omit(clause.Associations).Updates(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::UPDATES Failed to update account_event")
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (r *AccountEventRepository) FindByAccountAndEventId(accountId, eventId uuid.UUID, accountEvent *model.AccountEvent) error {
	if err := r.db.Where("account_id = ? AND event_id = ?", accountId, eventId).Preload("Account").First(&accountEvent).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_ACCOUNT_AND_EVENT_ID Failed to find account_event")
		}
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (r *AccountEventRepository) FindAccountsByEventId(eventId uuid.UUID, accounts *[]model.Account) error {
	if err := r.db.
		Table("account").
		Joins("JOIN account_event ON account.id = account_event.account_id").
		Where("account_event.event_id = ?", eventId).
		Find(accounts).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_ACCOUNTS_BY_EVENT_ID Failed to find accounts by event_id")
		return err
	}

	return nil
}
