package repository

import (
	"app/db"
	model "app/db/models"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AvailabilityRepository struct{}

// FindOverlappingAvailabilities finds overlapping availabilities for a given availability
func (*AvailabilityRepository) FindOverlappingAvailabilities(availability *model.Availability, availabilities *[]model.Availability) error {
	if err := db.GetDB().Where("account_id = ? AND event_id = ? AND starts_at <= ? AND ends_at >= ?",
		availability.AccountId,
		availability.EventId,
		availability.EndsAt,
		availability.StartsAt,
	).Find(&availabilities).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::FIND_OVERLAPPING_AVAILABILITIES Failed to find overlapping availabilities")
		return err
	}

	return nil
}

// DeleteByIds deletes availabilities by IDs
func (*AvailabilityRepository) DeleteByIds(ids *[]uuid.UUID) error {
	if err := db.GetDB().Where("id IN ?", *ids).Delete(&model.Availability{}).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_IDS Failed to delete availabilities")
		return err
	}

	return nil
}

// Create creates an availability
func (*AvailabilityRepository) Create(availability *model.Availability) error {
	if err := db.GetDB().Preload("Account").Create(&availability).First(&availability).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::CREATE Failed to create availability")
		return err
	}

	availability.Sanitized()

	return nil
}

// FindOneById finds an availability by ID
func (*AvailabilityRepository) FindOneById(id uuid.UUID, availability *model.Availability) error {
	if id == uuid.Nil {
		return errors.New("id is nil UUID")
	}
	if availability == nil {
		return errors.New("availability pointer is nil")
	}

	if err := db.GetDB().Preload("Account").Preload("Event").Preload("Event.AccountEvents").First(&availability, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::FIND_ONE_BY_ID Failed to find availability by ID")
		return err
	}

	return nil
}

// DeleteById deletes an availability by ID
func (*AvailabilityRepository) DeleteById(availabilityId *uuid.UUID) error {
	if availabilityId == nil {
		return errors.New("availabilityId pointer is nil")
	}

	if err := db.GetDB().Delete(&model.Availability{}, availabilityId).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_ID Failed to delete availability")
		return err
	}

	return nil
}

// retrieves all availabilities for a given event ID
func (*AvailabilityRepository) FindByEventId(eventId uuid.UUID, availabilities *[]model.Availability) error {
	if err := db.GetDB().Where("event_id = ?", eventId).Find(&availabilities).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("AVAILABILITY_REPOSITORY::GET_BY_EVENT_ID Failed to get availabilities by event ID")
		return err
	}

	return nil
}
