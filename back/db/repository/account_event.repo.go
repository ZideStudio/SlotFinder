package repository

import (
	"app/commons/constants"
	"app/db"
	model "app/db/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AccountEventRepository struct{}

func (*AccountEventRepository) Create(accountEvent *model.AccountEvent) error {
	if err := db.GetDB().Preload("Account").Create(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::CREATE Failed to create account_event")
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (*AccountEventRepository) Updates(accountEvent *model.AccountEvent) error {
	if err := db.GetDB().Where("account_id = ? AND event_id = ?", accountEvent.AccountId, accountEvent.EventId).Updates(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::UPDATES Failed to update account_event")
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (*AccountEventRepository) FindByAccountAndEventId(accountId, eventId uuid.UUID, accountEvent *model.AccountEvent) error {
	if err := db.GetDB().Where("account_id = ? AND event_id = ?", accountId, eventId).Preload("Account").First(&accountEvent).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_ACCOUNT_AND_EVENT_ID Failed to find account_event")
		return err
	}

	accountEvent.Sanitized()

	return nil
}

func (r *AccountEventRepository) FindEventIdsByAccountId(accountId uuid.UUID, limit int, offset int) ([]uuid.UUID, int64, error) {
	db := db.GetDB()

	var total int64
	if err := db.
		Model(&model.AccountEvent{}).
		Where("account_id = ?", accountId).
		Distinct("event_id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var eventIds []uuid.UUID
	if err := db.
		Model(&model.AccountEvent{}).
		Select("account_event.event_id").
		Joins("JOIN event ON account_event.event_id = event.id").
		Where("account_event.account_id = ?", accountId).
		Order(fmt.Sprintf("CASE WHEN event.status = '%s' THEN 1 WHEN event.status = '%s' THEN 2 WHEN event.status = '%s' THEN 3 ELSE 4 END, event.name ASC", constants.EVENT_STATUS_IN_DECISION, constants.EVENT_STATUS_UPCOMING, constants.EVENT_STATUS_FINISHED)).
		Limit(limit).
		Offset(offset).
		Pluck("account_event.event_id", &eventIds).Error; err != nil {
		return nil, 0, err
	}

	return eventIds, total, nil
}

func (r *AccountEventRepository) FindByIds(eventIds []uuid.UUID, accountEvents *[]model.AccountEvent) error {
	if err := db.GetDB().Preload("Account").Preload("Event").Preload("Event.Owner").Preload("Event.Availabilities").Preload("Event.Availabilities.Account").Preload("Event.Slots").Where("event_id IN ?", eventIds).Find(&accountEvents).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_IDS Failed to find account_events")
		return err
	}

	r.sanitizeAccountEvents(accountEvents)

	return nil
}

func (r *AccountEventRepository) FindByEventId(eventId uuid.UUID, accountEvents *[]model.AccountEvent) error {
	if err := db.GetDB().Where("event_id = ?", eventId).Preload("Account").Find(&accountEvents).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_EVENT_REPOSITORY::FIND_BY_EVENT_ID Failed to find account_events")
		return err
	}

	r.sanitizeAccountEvents(accountEvents)

	return nil
}

func (*AccountEventRepository) sanitizeAccountEvents(accountEvents *[]model.AccountEvent) {
	for i := range *accountEvents {
		(*accountEvents)[i].Sanitized()
	}
}
