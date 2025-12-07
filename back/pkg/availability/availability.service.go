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

// validateAvailabilityTimes validates the start and end times of an availability
func (s *AvailabilityService) validateAvailabilityTimes(startsAt, endsAt time.Time, event *model.Event) error {
	startsAt = startsAt.Truncate(time.Minute)
	endsAt = endsAt.Truncate(time.Minute)

	// Prevent creating/updating availabilities with end date before start date
	if startsAt.After(endsAt) {
		return constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent creating/updating availabilities with less than minimum duration of 5 minutes
	minDuration := 5 * time.Minute
	duration := endsAt.Sub(startsAt)
	if duration < minDuration {
		return constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err
	}

	// Prevent creating/updating availabilities not aligned on 5 minutes interval
	if startsAt.Minute()%5 != 0 || endsAt.Minute()%5 != 0 {
		return constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err
	}

	// Prevent creating/updating availabilities outside of event date range
	if startsAt.Before(event.StartsAt) {
		return constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err
	}
	if endsAt.After(event.EndsAt) {
		return constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err
	}

	return nil
}

// mergeAndCreateAvailability merges overlapping availabilities and creates the final availability
func (s *AvailabilityService) mergeAndCreateAvailability(tx *gorm.DB, availability *model.Availability, idsToDeleteFirst *[]uuid.UUID) (model.Availability, error) {
	var finalAvailability model.Availability

	// Delete any existing availabilities if specified (for update operation)
	if idsToDeleteFirst != nil && len(*idsToDeleteFirst) > 0 {
		if err := s.availabilityRepository.DeleteByIdsWithTx(tx, idsToDeleteFirst); err != nil {
			return model.Availability{}, err
		}
	}

	// Find overlapping availabilities within the locked transaction
	var availabilitiesToMerge []model.Availability
	if err := s.availabilityRepository.FindOverlappingAvailabilitiesWithTx(tx, availability, &availabilitiesToMerge); err != nil {
		return model.Availability{}, err
	}

	if len(availabilitiesToMerge) == 0 {
		// No overlapping availabilities, just create the new one
		if err := s.availabilityRepository.CreateWithTx(tx, availability); err != nil {
			return model.Availability{}, err
		}
		return *availability, nil
	}

	// Merge overlapping availabilities
	var availabilitiesIdsToDelete []uuid.UUID
	for _, existingAvailability := range availabilitiesToMerge {
		if existingAvailability.StartsAt.Before(availability.StartsAt) {
			availability.StartsAt = existingAvailability.StartsAt
		}
		if existingAvailability.EndsAt.After(availability.EndsAt) {
			availability.EndsAt = existingAvailability.EndsAt
		}

		availabilitiesIdsToDelete = append(availabilitiesIdsToDelete, existingAvailability.Id)
	}

	// Delete merged availabilities within the transaction
	if err := s.availabilityRepository.DeleteByIdsWithTx(tx, &availabilitiesIdsToDelete); err != nil {
		return model.Availability{}, err
	}

	// Create the merged availability within the transaction
	if err := s.availabilityRepository.CreateWithTx(tx, availability); err != nil {
		return model.Availability{}, err
	}

	finalAvailability = *availability
	return finalAvailability, nil
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

	// Validate availability times
	if err := s.validateAvailabilityTimes(data.StartsAt, data.EndsAt, &event); err != nil {
		return model.Availability{}, err
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
		var err error
		finalAvailability, err = s.mergeAndCreateAvailability(tx, &availabilityToCreate, nil)
		return err
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

	// Validate availability times
	if err := s.validateAvailabilityTimes(data.StartsAt, data.EndsAt, &availability.Event); err != nil {
		return model.Availability{}, err
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
		idsToDelete := []uuid.UUID{availabilityId}
		var err error
		finalAvailability, err = s.mergeAndCreateAvailability(tx, &availabilityToUpdate, &idsToDelete)
		return err
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
