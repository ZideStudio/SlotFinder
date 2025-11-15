package repository

import (
	"app/db"
	model "app/db/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AvailabilityRepository struct{}

func (*AvailabilityRepository) Create(availability *model.Availability) error {
	if err := db.GetDB().Preload("Account").Create(&availability).First(&availability).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::CREATE Failed to create availability")
		return err
	}

	availability.Sanitized()

	return nil
}

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

func (*AvailabilityRepository) DeleteByIds(ids *[]uuid.UUID) error {
	idsStrings := make([]string, len(*ids))
	for i, id := range *ids {
		idsStrings[i] = id.String()
	}

	if err := db.GetDB().Where("id IN (?)", idsStrings).Delete(&model.Availability{}).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_IDS Failed to delete availabilities")
		return err
	}

	return nil
}
