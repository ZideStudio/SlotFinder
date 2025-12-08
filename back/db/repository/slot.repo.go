package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SlotRepository struct{}

func (*SlotRepository) Create(slot *model.Slot) error {
	if err := db.GetDB().Create(&slot).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::CREATE Failed to create slot")
		return err
	}

	return nil
}

func (*SlotRepository) Updates(slot *model.Slot) error {
	if err := db.GetDB().Updates(&slot).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::UPDATES Failed to update slot")
		return err
	}

	return nil
}

func (*SlotRepository) FindOneById(slotId uuid.UUID, slot *model.Slot) error {
	if err := db.GetDB().Where("id = ?", slotId.String()).Preload("Event").Preload("Event.AccountEvents").First(slot).Error; err != nil {
		log.Error().Err(err).Str("slotId", slotId.String()).Msg("SLOT_REPOSITORY::FIND_ONE_BY_ID Failed to find slot by id")
		return err
	}

	return nil
}

func (*SlotRepository) FindByEventId(eventId uuid.UUID, slots *[]model.Slot) error {
	if err := db.GetDB().Where("event_id = ?", eventId).Find(slots).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::FIND_BY_EVENT_ID Failed to find slots by event id")
		return err
	}

	return nil
}

func (*SlotRepository) DeleteByEventId(eventId uuid.UUID) error {
	if err := db.GetDB().Where("event_id = ?", eventId).Delete(&model.Slot{}).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::DELETE_BY_EVENT_ID Failed to delete slots by event id")
		return err
	}

	return nil
}

func (*SlotRepository) Delete(id uuid.UUID) error {
	if err := db.GetDB().Where("id = ?", id.String()).Delete(&model.Slot{}).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::DELETE Failed to delete slot")
		return err
	}

	return nil
}

func (*SlotRepository) DeleteByIds(ids []uuid.UUID) error {
	if err := db.GetDB().Where("id IN ?", ids).Delete(&model.Slot{}).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::DELETE_BY_IDS Failed to delete slots by ids")
		return err
	}

	return nil
}
