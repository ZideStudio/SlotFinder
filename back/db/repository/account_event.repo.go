package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AccountEventRepository struct{}

func (*AccountEventRepository) Create(accountEvent *model.AccountEvent) error {
	if err := db.GetDB().Preload("Account").Create(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::CREATE Failed to create account_event")
		return err
	}

	accountEvent = accountEvent.Sanitized()

	return nil
}

func (*AccountEventRepository) FindByAccountAndEventId(accountId, eventId uuid.UUID, accountEvent *model.AccountEvent) error {
	if err := db.GetDB().Where("account_id = ? AND event_id = ?", accountId, eventId).Preload("Account").First(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_ACCOUNT_AND_EVENT_ID Failed to find account_event")
		return err
	}

	accountEvent = accountEvent.Sanitized()

	return nil
}

func (r *AccountEventRepository) FindByAccountId(accountId uuid.UUID, accountEvents *[]model.AccountEvent) error {
	if err := db.GetDB().Where("account_id = ?", accountId).Preload("Account").Preload("Event").Preload("Event.Owner").Find(&accountEvents).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_ACCOUNT_ID Failed to find account_events")
		return err
	}

	r.sanitizeAccountEvents(accountEvents)

	return nil
}

func (r *AccountEventRepository) FindByIds(eventIds []uuid.UUID, accountEvents *[]model.AccountEvent) error {
	if err := db.GetDB().Preload("Account").Preload("Event").Preload("Event.Owner").Preload("Event.Availabilities").Preload("Event.Availabilities.Account").Where("event_id IN ?", eventIds).Find(&accountEvents).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_IDS Failed to find account_events")
		return err
	}

	r.sanitizeAccountEvents(accountEvents)

	return nil
}

func (*AccountEventRepository) sanitizeAccountEvents(accountEvents *[]model.AccountEvent) {
	for i := range *accountEvents {
		(*accountEvents)[i] = *(*accountEvents)[i].Sanitized()
	}
}
