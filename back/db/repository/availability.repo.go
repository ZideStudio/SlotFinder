package repository

import (
	"app/db"
	model "app/db/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AvailabilityRepository struct{}

// CreateWithLock creates an availability with row-level locking to prevent concurrent duplicates
func (*AvailabilityRepository) CreateWithLock(availability *model.Availability, mergeFunc func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock existing availabilities for this user and event, this prevents concurrent requests from creating duplicates
		var lockRows []model.Availability
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where(
			"account_id = ? AND event_id = ?",
			availability.AccountId,
			availability.EventId,
		).Find(&lockRows).Error; err != nil {
			log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::CREATE_WITH_LOCK Failed to lock rows")
			return err
		}

		// Execute the merge function with the locked transaction
		if err := mergeFunc(tx); err != nil {
			return err
		}

		return nil
	})
}

// FindOverlappingAvailabilitiesWithTx finds overlapping availabilities within a transaction
func (*AvailabilityRepository) FindOverlappingAvailabilitiesWithTx(tx *gorm.DB, availability *model.Availability, availabilities *[]model.Availability) error {
	if err := tx.Where("account_id = ? AND event_id = ? AND starts_at <= ? AND ends_at >= ?",
		availability.AccountId,
		availability.EventId,
		availability.EndsAt,
		availability.StartsAt,
	).Find(&availabilities).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::FIND_OVERLAPPING_AVAILABILITIES_WITH_TX Failed to find overlapping availabilities")
		return err
	}

	return nil
}

// DeleteByIdsWithTx deletes availabilities by IDs within a transaction
func (*AvailabilityRepository) DeleteByIdsWithTx(tx *gorm.DB, ids *[]uuid.UUID) error {
	idsStrings := make([]string, len(*ids))
	for i, id := range *ids {
		idsStrings[i] = id.String()
	}

	if err := tx.Where("id IN (?)", idsStrings).Delete(&model.Availability{}).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_IDS_WITH_TX Failed to delete availabilities")
		return err
	}

	return nil
}

// CreateWithTx creates an availability within a transaction
func (*AvailabilityRepository) CreateWithTx(tx *gorm.DB, availability *model.Availability) error {
	if err := tx.Preload("Account").Create(&availability).First(&availability).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::CREATE_WITH_TX Failed to create availability")
		return err
	}

	availability.Sanitized()

	return nil
}
