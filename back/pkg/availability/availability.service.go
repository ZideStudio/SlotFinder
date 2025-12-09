package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/slot"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilityService struct {
	slotService            *slot.SlotService
	availabilityRepository *repository.AvailabilityRepository
	eventRepository        *repository.EventRepository
	locks                  sync.Map // Map to store mutexes per (accountId, eventId) combination
}

func NewAvailabilityService(service *AvailabilityService) *AvailabilityService {
	if service != nil {
		return service
	}

	return &AvailabilityService{
		slotService:            slot.NewSlotService(nil),
		availabilityRepository: &repository.AvailabilityRepository{},
		eventRepository:        &repository.EventRepository{},
	}
}

// validateAvailabilityTimes validates the time constraints for an availability
func (s *AvailabilityService) validateAvailabilityTimes(startsAt, endsAt time.Time, event *model.Event) error {
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
	// Check if times are exactly on 5-minute boundaries (no seconds or sub-seconds)
	if startsAt.Truncate(5*time.Minute) != startsAt || endsAt.Truncate(5*time.Minute) != endsAt {
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

// validateEventAccess validates that the event exists, is accessible, and not ended
func (s *AvailabilityService) validateEventAccess(eventId uuid.UUID, userId *uuid.UUID, event *model.Event) error {
	if err := s.eventRepository.FindOneById(eventId, event); err != nil {
		return constants.ERR_EVENT_NOT_FOUND.Err
	}

	// Check if event is ended
	if event.IsLocked() {
		return constants.ERR_EVENT_ENDED.Err
	}

	// Check if user has access to the event
	if !event.HasUserAccess(userId) {
		return constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	return nil
}

func (s *AvailabilityService) Create(data *AvailabilityCreateDto, eventId uuid.UUID, user *guard.Claims) (model.Availability, error) {
	// Acquire per-event mutex to prevent concurrent availability modifications
	value, _ := s.locks.LoadOrStore(user.Id.String(), &sync.Mutex{})
	mu := value.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	// Get event and validate access
	var event model.Event
	if err := s.validateEventAccess(eventId, &user.Id, &event); err != nil {
		return model.Availability{}, err
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

	// Find overlapping availabilities
	var availabilitiesToMerge []model.Availability
	if err := s.availabilityRepository.FindOverlappingAvailabilities(&availabilityToCreate, &availabilitiesToMerge); err != nil {
		return model.Availability{}, err
	}

	if len(availabilitiesToMerge) == 0 {
		// No overlapping availabilities, just create the new one
		if err := s.availabilityRepository.Create(&availabilityToCreate); err != nil {
			return model.Availability{}, err
		}
		return availabilityToCreate, nil
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

	// Delete merged availabilities
	if err := s.availabilityRepository.DeleteByIds(&availabilitiesIdsToDelete); err != nil {
		return model.Availability{}, err
	}

	// Create the merged availability
	if err := s.availabilityRepository.Create(&availabilityToCreate); err != nil {
		return model.Availability{}, err
	}

	// Trigger slot recalculation asynchronously
	go s.slotService.LoadSlots(eventId)

	return availabilityToCreate, nil
}

func (s *AvailabilityService) Update(data *AvailabilityUpdateDto, availabilityId uuid.UUID, user *guard.Claims) (model.Availability, error) {
	// Acquire per-event mutex to prevent concurrent availability modifications
	value, _ := s.locks.LoadOrStore(user.Id.String(), &sync.Mutex{})
	mu := value.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

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

	// Check if event is locked (same check as in Create via validateEventAccess)
	if availability.Event.IsLocked() {
		return model.Availability{}, constants.ERR_EVENT_ENDED.Err
	}

	// Update fields if provided
	updated := false
	if data.StartsAt != nil {
		availability.StartsAt = data.StartsAt.Truncate(time.Minute)
		updated = true
	}
	if data.EndsAt != nil {
		availability.EndsAt = data.EndsAt.Truncate(time.Minute)
		updated = true
	}

	// If no fields were updated, return early
	if !updated {
		return availability, nil
	}

	// Validate availability times
	if err := s.validateAvailabilityTimes(availability.StartsAt, availability.EndsAt, &availability.Event); err != nil {
		return model.Availability{}, err
	}

	// Update availability
	if err := s.availabilityRepository.Update(&availability); err != nil {
		return model.Availability{}, err
	}

	// Trigger slot recalculation asynchronously
	go s.slotService.LoadSlots(availability.EventId)

	return availability, nil
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

	// Trigger slot recalculation asynchronously
	go s.slotService.LoadSlots(availability.EventId)

	return nil
}
