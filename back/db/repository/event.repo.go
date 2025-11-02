package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EventRepository struct{}

func (*EventRepository) Create(event *model.Event) error {
	if err := db.GetDB().Create(&event).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::CREATE Failed to create event")
		return err
	}

	return nil
}

func (*EventRepository) FindOneById(id uuid.UUID, event *model.Event) error {
	if err := db.GetDB().Where("id = ?", id.String()).Preload("Owner").Preload("AccountEvents").Preload("AccountEvents.Account").First(&event).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::FIND_ONE_BY_ID Failed to find event by id")
		return err
	}

	return nil
}

func (*EventRepository) Delete(id uuid.UUID) error {
	if err := db.GetDB().Where("id = ?", id.String()).Delete(&model.Event{}).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::DELETE Failed to delete event")
		return err
	}

	return nil
}
