package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"app/pkg/slot"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService struct {
	eventRepository        *repository.EventRepository
	accountEventRepository *repository.AccountEventRepository
	availabilityRepository *repository.AvailabilityRepository
	slotRepository         *repository.SlotRepository
	slotService            *slot.SlotService
	signinService          *signin.SigninService
}

func NewEventService(service *EventService) *EventService {
	if service != nil {
		return service
	}

	return &EventService{
		eventRepository:        &repository.EventRepository{},
		accountEventRepository: &repository.AccountEventRepository{},
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		slotRepository:         &repository.SlotRepository{},
		slotService:            slot.NewSlotService(nil),
		signinService:          signin.NewSigninService(nil),
	}
}

func (s *EventService) Create(data *EventCreateDto, user *guard.Claims) (model.Event, error) {
	// Prevent creating events with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return model.Event{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent creating events in the past
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if data.StartsAt.Before(now) {
		return model.Event{}, constants.ERR_EVENT_START_BEFORE_TODAY.Err
	}

	// Prevent creating events with duration less than 1 day
	oneDayAfterStart := data.StartsAt.Add(24 * time.Hour)
	if data.EndsAt.Before(oneDayAfterStart) {
		return model.Event{}, constants.ERR_EVENT_DURATION_TOO_SHORT.Err
	}

	// Create event
	event := model.Event{
		Id:          uuid.New(),
		Name:        data.Name,
		Description: data.Description,
		Duration:    data.Duration,
		StartsAt:    data.StartsAt,
		EndsAt:      data.EndsAt,
		Owner: model.Account{
			Id:       user.Id,
			UserName: user.Username,
		},
		Status: constants.EVENT_STATUS_IN_DECISION,
	}
	if err := s.eventRepository.Create(&event); err != nil {
		return event, err
	}

	// Create account_event relation
	accountEvent := model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		_ = s.eventRepository.Delete(event.Id)
		return event, err
	}
	if err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent); err != nil {
		_ = s.eventRepository.Delete(event.Id)
		return event, err
	}

	return s.GetEvent(event.Id, user)
}

// SetEventDatesFromDto validates and sets the event dates from the provided DTO values.
func SetEventDatesFromDto(event *model.Event, startsAtDto, endsAtDto *time.Time) error {
	if event == nil {
		return errors.New("event pointer is nil")
	}
	if startsAtDto == nil && endsAtDto == nil {
		return nil
	}

	hasStartDate := startsAtDto != nil
	startsAt := event.StartsAt
	endsAt := event.EndsAt
	if hasStartDate {
		startsAt = *startsAtDto
	}
	if endsAtDto != nil {
		endsAt = *endsAtDto
	}

	// Prevent creating events with end date before start date
	if startsAt.After(endsAt) {
		return constants.ERR_EVENT_START_AFTER_END.Err
	}
	// Prevent creating events with duration less than 1 day
	oneDayAfterStart := startsAt.Add(24 * time.Hour)
	if endsAt.Before(oneDayAfterStart) {
		return constants.ERR_EVENT_DURATION_TOO_SHORT.Err
	}
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if hasStartDate && startsAt.Before(now) { // Prevent updating events start date before today
		return constants.ERR_EVENT_START_BEFORE_TODAY.Err
	} else if endsAt.Before(now) { // Prevent updating events end date before today
		return constants.ERR_EVENT_START_BEFORE_TODAY.Err
	}

	if validatedSlot := event.GetValidatedSlot(); validatedSlot != nil && (endsAt.Before(validatedSlot.EndsAt) || startsAt.After(validatedSlot.StartsAt)) {
		// Prevent updating events to end date before already validated slot
		return constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err
	}

	// Set parsed dates back to event
	event.StartsAt = startsAt
	event.EndsAt = endsAt

	return nil
}

func (s *EventService) Update(eventId uuid.UUID, data *EventUpdateDto, user *guard.Claims) error {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ERR_EVENT_NOT_FOUND.Err
		}
		return err
	}

	// Check if user is the owner of the event
	if !event.IsOwner(&user.Id) {
		return constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	// Data validation
	if data.Status != nil && *data.Status != constants.EVENT_STATUS_IN_DECISION {
		return errors.New("only status 'in_decision' can be set manually")
	}
	var isBreakingSlots bool
	if data.StartsAt != nil || data.EndsAt != nil {
		if err := SetEventDatesFromDto(&event, data.StartsAt, data.EndsAt); err != nil {
			return err
		}
		isBreakingSlots = true
	}

	// Update fields if provided
	if data.Name != nil {
		event.Name = *data.Name
	}
	if data.Description != nil {
		event.Description = data.Description
	}
	if data.Duration != nil {
		event.Duration = *data.Duration
		isBreakingSlots = true
	}
	var isStatusChanged bool
	if data.Status != nil {
		event.Status = *data.Status
		isStatusChanged = true
	}

	// Update event in repository
	if err := s.eventRepository.Updates(&event); err != nil {
		return err
	}

	// If status changed, remove validated slot
	if isStatusChanged {
		err := s.slotRepository.DeleteValidatedSlotByEventId(event.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// If dates are not being updated, return
	if !isBreakingSlots {
		return nil
	}

	// Remove slots
	if err := s.slotRepository.DeleteByEventId(event.Id); err != nil {
		return err
	}

	// Remove availabilities that are out of the new event date range
	if err := s.availabilityRepository.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt); err != nil {
		return err
	}

	// Load slots
	go s.slotService.LoadSlots(eventId)

	return nil
}

func (s *EventService) GetUserEvents(
	user *guard.Claims,
	pagination *lib.Pagination[model.Event],
) error {
	events, total, err := s.eventRepository.FindEventsByAccountId(user.Id, pagination.Limit, pagination.Offset)
	if err != nil {
		return err
	}
	pagination.Total = total
	pagination.Data = events

	return nil
}

func (s *EventService) GetEvent(eventId uuid.UUID, user *guard.Claims) (model.Event, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, constants.ERR_EVENT_NOT_FOUND.Err
		}

		return event, err
	}

	// If no user, return event basic info
	if user == nil {
		event.Participants = []model.Account{}
		return event, nil
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	notJoined := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !notJoined {
		return event, err
	}
	if notJoined {
		event.Participants = []model.Account{}
		return event, nil
	}

	// Return event info
	return event, nil
}

func (s *EventService) JoinEvent(eventId uuid.UUID, user *guard.Claims) (model.Event, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, constants.ERR_EVENT_NOT_FOUND.Err
		}

		return event, err
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	alreadyJoined := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return event, err
	}
	if alreadyJoined {
		return event, constants.ERR_EVENT_ALREADY_JOINED.Err
	}

	// Create account_event relation
	accountEvent = model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		return event, err
	}

	return s.GetEvent(event.Id, user)
}

func (s *EventService) UpdateProfile(data *EventProfileDto, eventId uuid.UUID, user *guard.Claims) error {
	// Find account event relation
	var accountEvent model.AccountEvent
	if err := s.accountEventRepository.FindByAccountAndEventId(user.Id, eventId, &accountEvent); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ERR_EVENT_NOT_FOUND.Err
		}
		return err
	}

	if !lib.IsHexa(data.Color) {
		return constants.ERR_INVALID_COLOR_FORMAT.Err
	}

	// Update color property
	accountEvent.Color = &data.Color
	if err := s.accountEventRepository.Updates(&accountEvent); err != nil {
		return err
	}

	return nil
}
