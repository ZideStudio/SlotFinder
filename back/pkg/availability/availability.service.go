package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilityService struct {
	availabilityRepository *repository.AvailabilityRepository
	eventRepository        *repository.EventRepository
}

func NewAvailabilityService(service *AvailabilityService) *AvailabilityService {
	if service != nil {
		return service
	}

	return &AvailabilityService{
		availabilityRepository: &repository.AvailabilityRepository{},
		eventRepository:        &repository.EventRepository{},
	}
}

func (s *AvailabilityService) Create(data *AvailabilityCreateDto, eventId uuid.UUID, user *guard.Claims) (model.Availability, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		return model.Availability{}, constants.ERR_EVENT_NOT_FOUND.Err
	}

	// Check if event is ended
	if event.HasEnded() {
		return model.Availability{}, constants.ERR_EVENT_ENDED.Err
	}

	// Check if user has access to the event
	if !event.HasUserAccess(&user.Id) {
		return model.Availability{}, constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	data.StartsAt = data.StartsAt.Truncate(time.Minute)
	data.EndsAt = data.EndsAt.Truncate(time.Minute)

	// Prevent creating availabilities with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return model.Availability{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent creating availabilities with less than minimum duration of 5 minutes
	minDuration := 5 * time.Minute
	duration := data.EndsAt.Sub(data.StartsAt)
	if duration < minDuration {
		return model.Availability{}, constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err
	}

	// Prevent creating availabilities not aligned on 5 minutes interval
	if data.StartsAt.Minute()%5 != 0 || data.EndsAt.Minute()%5 != 0 {
		return model.Availability{}, constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err
	}

	// Prevent creating availabilities outside of event date range
	if data.StartsAt.Before(event.StartsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err
	}
	if data.EndsAt.After(event.EndsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err
	}

	// Create availability model
	availabilityToCreate := model.Availability{
		Id:        uuid.New(),
		StartsAt:  data.StartsAt,
		EndsAt:    data.EndsAt,
		AccountId: user.Id,
		EventId:   eventId,
	}

	// Use transaction with row-level locking to prevent concurrent duplicates
	var finalAvailability model.Availability
	err := s.availabilityRepository.CreateWithLock(&availabilityToCreate, func(tx *gorm.DB) error {
		// Find overlapping availabilities within the locked transaction
		var availabilitiesToMerge []model.Availability
		if err := s.availabilityRepository.FindOverlappingAvailabilitiesWithTx(tx, &availabilityToCreate, &availabilitiesToMerge); err != nil {
			return err
		}

		if len(availabilitiesToMerge) == 0 {
			// No overlapping availabilities, just create the new one
			if err := s.availabilityRepository.CreateWithTx(tx, &availabilityToCreate); err != nil {
				return err
			}
			finalAvailability = availabilityToCreate
			return nil
		}

		// Merge overlapping availabilities
		var availabilitiesIdsToDelete []uuid.UUID
		for _, existingAvailability := range availabilitiesToMerge {
			if existingAvailability.StartsAt.Before(availabilityToCreate.StartsAt) {
				availabilityToCreate.StartsAt = existingAvailability.StartsAt
			}
			if existingAvailability.EndsAt.After(availabilityToCreate.EndsAt) {
				availabilityToCreate.EndsAt = existingAvailability.EndsAt
			}

			availabilitiesIdsToDelete = append(availabilitiesIdsToDelete, existingAvailability.Id)
		}

		// Delete merged availabilities within the transaction
		if err := s.availabilityRepository.DeleteByIdsWithTx(tx, &availabilitiesIdsToDelete); err != nil {
			return err
		}

		// Create the merged availability within the transaction
		if err := s.availabilityRepository.CreateWithTx(tx, &availabilityToCreate); err != nil {
			return err
		}

		finalAvailability = availabilityToCreate
		return nil
	})

	if err != nil {
		return model.Availability{}, err
	}

	return finalAvailability, nil
}

func (s *AvailabilityService) Update(data *AvailabilityUpdateDto, availabilityId uuid.UUID, user *guard.Claims) (model.Availability, error) {
	// Get availability
	var availability model.Availability
	if err := s.availabilityRepository.FindOneById(availabilityId, &availability); err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Availability{}, constants.ERR_AVAILABILITY_NOT_FOUND.Err
		}
		return model.Availability{}, err
	}

	// Check if availability belongs to the user
	if availability.AccountId != user.Id {
		return model.Availability{}, constants.ERR_AVAILABILITY_ACCESS_DENIED.Err
	}

	// Check if user has access to the event
	if !availability.Event.HasUserAccess(&user.Id) {
		return model.Availability{}, constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	// Check if event is ended
	if availability.Event.HasEnded() {
		return model.Availability{}, constants.ERR_EVENT_ENDED.Err
	}

	data.StartsAt = data.StartsAt.Truncate(time.Minute)
	data.EndsAt = data.EndsAt.Truncate(time.Minute)

	// Prevent updating availabilities with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return model.Availability{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent updating availabilities with less than minimum duration of 5 minutes
	minDuration := 5 * time.Minute
	duration := data.EndsAt.Sub(data.StartsAt)
	if duration < minDuration {
		return model.Availability{}, constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err
	}

	// Prevent updating availabilities not aligned on 5 minutes interval
	if data.StartsAt.Minute()%5 != 0 || data.EndsAt.Minute()%5 != 0 {
		return model.Availability{}, constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err
	}

	// Prevent updating availabilities outside of event date range
	if data.StartsAt.Before(availability.Event.StartsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err
	}
	if data.EndsAt.After(availability.Event.EndsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err
	}

	// Update availability model
	availabilityToUpdate := model.Availability{
		Id:        availabilityId,
		StartsAt:  data.StartsAt,
		EndsAt:    data.EndsAt,
		AccountId: user.Id,
		EventId:   availability.EventId,
	}

	// Use transaction with row-level locking to prevent concurrent duplicates
	var finalAvailability model.Availability
	err := s.availabilityRepository.UpdateWithLock(&availabilityToUpdate, func(tx *gorm.DB) error {
		// Delete the old availability first
		if err := s.availabilityRepository.DeleteByIdsWithTx(tx, &[]uuid.UUID{availabilityId}); err != nil {
			return err
		}

		// Find overlapping availabilities within the locked transaction (excluding the one being updated)
		var availabilitiesToMerge []model.Availability
		if err := s.availabilityRepository.FindOverlappingAvailabilitiesWithTx(tx, &availabilityToUpdate, &availabilitiesToMerge); err != nil {
			return err
		}

		if len(availabilitiesToMerge) == 0 {
			// No overlapping availabilities, just create the updated one
			if err := s.availabilityRepository.CreateWithTx(tx, &availabilityToUpdate); err != nil {
				return err
			}
			finalAvailability = availabilityToUpdate
			return nil
		}

		// Merge overlapping availabilities
		var availabilitiesIdsToDelete []uuid.UUID
		for _, existingAvailability := range availabilitiesToMerge {
			if existingAvailability.StartsAt.Before(availabilityToUpdate.StartsAt) {
				availabilityToUpdate.StartsAt = existingAvailability.StartsAt
			}
			if existingAvailability.EndsAt.After(availabilityToUpdate.EndsAt) {
				availabilityToUpdate.EndsAt = existingAvailability.EndsAt
			}

			availabilitiesIdsToDelete = append(availabilitiesIdsToDelete, existingAvailability.Id)
		}

		// Delete merged availabilities within the transaction
		if err := s.availabilityRepository.DeleteByIdsWithTx(tx, &availabilitiesIdsToDelete); err != nil {
			return err
		}

		// Create the merged availability within the transaction
		if err := s.availabilityRepository.CreateWithTx(tx, &availabilityToUpdate); err != nil {
			return err
		}

		finalAvailability = availabilityToUpdate
		return nil
	})

	if err != nil {
		return model.Availability{}, err
	}

	return finalAvailability, nil
}

func (s *AvailabilityService) Delete(availabilityId uuid.UUID, user *guard.Claims) error {
	// Get availability
	var availability model.Availability
	if err := s.availabilityRepository.FindOneById(availabilityId, &availability); err != nil {
		if err == gorm.ErrRecordNotFound {
			return constants.ERR_AVAILABILITY_NOT_FOUND.Err
		}
		return err
	}

	// Check if availability belongs to the user
	if availability.AccountId != user.Id {
		return constants.ERR_AVAILABILITY_ACCESS_DENIED.Err
	}

	// Check if user has access to the event
	if !availability.Event.HasUserAccess(&user.Id) {
		return constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	// Check if event is ended
	if availability.Event.HasEnded() {
		return constants.ERR_EVENT_ENDED.Err
	}

	// Delete availability
	if err := s.availabilityRepository.DeleteById(&availabilityId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return constants.ERR_AVAILABILITY_NOT_FOUND.Err
		}
		return err
	}

	return nil
}
