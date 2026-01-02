package repository

import (
	"app/db"
	model "app/db/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AvailabilityRepository struct {
	db *gorm.DB
}

func NewAvailabilityRepository(database *gorm.DB) *AvailabilityRepository {
	if database == nil {
		database = db.GetDB()
	}
	return &AvailabilityRepository{
		db: database,
	}
}

// Finds overlapping availabilities for a given availability
func (r *AvailabilityRepository) FindOverlappingAvailabilities(availability *model.Availability, availabilities *[]model.Availability) error {
	if err := r.db.Where("account_id = ? AND event_id = ? AND starts_at <= ? AND ends_at >= ?",
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

// Deletes availabilities by IDs
func (r *AvailabilityRepository) DeleteByIds(ids *[]uuid.UUID) error {
	if err := r.db.Where("id IN ?", *ids).Delete(&model.Availability{}).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_IDS Failed to delete availabilities")
		return err
	}

	return nil
}

// Creates an availability
func (r *AvailabilityRepository) Create(availability *model.Availability) error {
	if err := r.db.Preload("Account").Create(&availability).First(&availability).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::CREATE Failed to create availability")
		return err
	}

	availability.Sanitized()

	return nil
}

// Finds an availability by ID
func (r *AvailabilityRepository) FindOneById(id uuid.UUID, availability *model.Availability) error {
	if id == uuid.Nil {
		return errors.New("id is nil UUID")
	}
	if availability == nil {
		return errors.New("availability pointer is nil")
	}

	if err := r.db.Preload("Account").Preload("Event").Preload("Event.AccountEvents").First(&availability, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::FIND_ONE_BY_ID Failed to find availability by ID")
		return err
	}

	return nil
}

// Deletes an availability by ID
func (r *AvailabilityRepository) DeleteById(availabilityId *uuid.UUID) error {
	if availabilityId == nil {
		return errors.New("availabilityId pointer is nil")
	}

	if err := r.db.Delete(&model.Availability{}, availabilityId).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_BY_ID Failed to delete availability")
		return err
	}

	return nil
}

// retrieves all availabilities for a given event ID
func (r *AvailabilityRepository) FindByEventId(eventId uuid.UUID, availabilities *[]model.Availability) error {
	if err := r.db.Where("event_id = ?", eventId).Find(&availabilities).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("AVAILABILITY_REPOSITORY::GET_BY_EVENT_ID Failed to get availabilities by event ID")
		return err
	}

	return nil
}

// Updates an availability
func (r *AvailabilityRepository) Update(availability *model.Availability) error {
	if availability == nil {
		return errors.New("availability pointer is nil")
	}

	if err := r.db.Model(&availability).Updates(availability).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::UPDATE Failed to update availability")
		return err
	}

	if err := r.db.Preload("Account").Preload("Event").Preload("Event.AccountEvents").First(&availability, "id = ?", availability.Id).Error; err != nil {
		log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::UPDATE Failed to reload availability after update")
		return err
	}

	availability.Sanitized()

	return nil
}

// DeleteOutOfEventRangeAndAdjustOverlaps deletes availabilities that are out of the event range and adjusts overlapping ones
func (r *AvailabilityRepository) DeleteOutOfEventRangeAndAdjustOverlaps(eventId uuid.UUID, startsAt time.Time, endsAt time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Find availabilities that overlap with the new event range and extend beyond it
		var overlappingAvailabilities []model.Availability

		// Build the WHERE clause with clear conditions
		overlapCondition := "(starts_at < ? AND ends_at > ?) AND (starts_at < ? OR ends_at > ?)"
		fullWhereClause := "event_id = ? AND " + overlapCondition

		if err := tx.Where(fullWhereClause,
			eventId,
			endsAt,
			startsAt,
			startsAt,
			endsAt,
		).Find(&overlappingAvailabilities).Error; err != nil {
			log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_OUT_OF_EVENT_RANGE_AND_ADJUST_OVERLAPS Failed to find overlapping availabilities")
			return err
		}

		// Adjust overlapping availabilities to fit within the event range
		for _, availability := range overlappingAvailabilities {
			adjustedAvailability := availability

			// Adjust start time if it's before the event starts
			if availability.StartsAt.Before(startsAt) {
				adjustedAvailability.StartsAt = startsAt
			}

			// Adjust end time if it's after the event ends
			if availability.EndsAt.After(endsAt) {
				adjustedAvailability.EndsAt = endsAt
			}

			// Update the availability with adjusted times
			if err := tx.Model(&availability).Updates(model.Availability{
				StartsAt: adjustedAvailability.StartsAt,
				EndsAt:   adjustedAvailability.EndsAt,
			}).Error; err != nil {
				log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_OUT_OF_EVENT_RANGE_AND_ADJUST_OVERLAPS Failed to update overlapping availability")
				return err
			}
		}

		// Delete availabilities that are completely out of range
		if err := tx.Where("event_id = ? AND (ends_at < ? OR starts_at > ?)", eventId, startsAt, endsAt).Delete(&model.Availability{}).Error; err != nil {
			log.Error().Err(err).Msg("AVAILABILITY_REPOSITORY::DELETE_OUT_OF_EVENT_RANGE_AND_ADJUST_OVERLAPS Failed to delete availabilities out of event range")
			return err
		}

		return nil
	})
}
