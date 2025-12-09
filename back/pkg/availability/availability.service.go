package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/slot"
	"fmt"
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

func (s *AvailabilityService) Create(data *AvailabilityCreateDto, eventId uuid.UUID, user *guard.Claims) (model.Availability, error) {
	// Acquire per-event mutex to prevent concurrent availability modifications
	lockKey := fmt.Sprintf("%s:%s", user.Id.String(), eventId.String())

	value, _ := s.locks.LoadOrStore(lockKey, &sync.Mutex{})
	mu := value.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		return model.Availability{}, constants.ERR_EVENT_NOT_FOUND.Err
	}

	// Check if event is ended
	if event.IsLocked() {
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

	// Acquire per-event mutex to prevent concurrent availability modifications
	lockKey := fmt.Sprintf("%s:%s", user.Id.String(), availability.EventId.String())

	value, _ := s.locks.LoadOrStore(lockKey, &sync.Mutex{})
	mu := value.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

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

	// Prevent updating availabilities with end date before start date
	if availability.StartsAt.After(availability.EndsAt) {
		return model.Availability{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent updating availabilities with less than minimum duration of 5 minutes
	minDuration := 5 * time.Minute
	duration := availability.EndsAt.Sub(availability.StartsAt)
	if duration < minDuration {
		return model.Availability{}, constants.ERR_AVAILABILITY_DURATION_TOO_SHORT.Err
	}

	// Prevent updating availabilities not aligned on 5 minutes interval
	if availability.StartsAt.Minute()%5 != 0 || availability.EndsAt.Minute()%5 != 0 {
		return model.Availability{}, constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL.Err
	}

	// Prevent updating availabilities outside of event date range
	if availability.StartsAt.Before(availability.Event.StartsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_START_BEFORE_EVENT.Err
	}
	if availability.EndsAt.After(availability.Event.EndsAt) {
		return model.Availability{}, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err
	}

	// Update availability
	if err := s.availabilityRepository.Update(&availability); err != nil {
		return model.Availability{}, err
	}

	// Trigger slot recalculation asynchronously
	go s.slotService.LoadSlots(availability.EventId)

	availability.Sanitized()

	return availability, nil
}
