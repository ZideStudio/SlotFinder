package repository

import (
	"app/db"
	model "app/db/models"
	"context"
	"errors"
	"hash/fnv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AvailabilityRepository struct{}

// computeAdvisoryLockKey computes a deterministic int64 key from account_id and event_id
// for use with PostgreSQL advisory locks. This ensures mutual exclusion across the entire
// operation, including both existing rows and potential new insertions.
func (*AvailabilityRepository) computeAdvisoryLockKey(accountId, eventId uuid.UUID) int64 {
	h := fnv.New64a()
	h.Write(accountId[:])
	h.Write(eventId[:])
	return int64(h.Sum64())
}

// CreateWithLock creates an availability with PostgreSQL advisory locking to prevent concurrent duplicates.
// Uses pg_advisory_xact_lock to ensure mutual exclusion for the entire operation, including both
// existing rows and potential new insertions. The lock is automatically released at transaction end.
func (r *AvailabilityRepository) CreateWithLock(availability *model.Availability, mergeFunc func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Acquire PostgreSQL advisory lock for this (account_id, event_id) combination
		// This prevents concurrent transactions from processing the same combination simultaneously
		lockKey := r.computeAdvisoryLockKey(availability.AccountId, availability.EventId)

		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			log.Error().Err(err).Int64("lockKey", lockKey).Msg("AVAILABILITY_REPOSITORY::CREATE_WITH_LOCK Failed to acquire advisory lock")
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
	if err := tx.Where("id IN ?", *ids).Delete(&model.Availability{}).Error; err != nil {
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

// Find an availability by ID and event id
func (*AvailabilityRepository) FindOneByIdAndEventId(id uuid.UUID, eventId uuid.UUID, availability *model.Availability) error {
	if id == uuid.Nil {
		return errors.New("id is nil UUID")
	}
	if availability == nil {
		return errors.New("availability pointer is nil")
	}

	if err := db.GetDB().Preload("Account").First(&availability, "id = ? AND event_id = ?", id, eventId).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::FIND_ONE_BY_ID Failed to find availability by ID")
		return err
	}

	return nil
}

// Deletes an availability by ID
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
