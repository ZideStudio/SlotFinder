package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/lib"
	"app/config"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
	"app/pkg/signin"
	"app/pkg/slot"
	"errors"
	"strings"
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
	mailService            *mail.MailService
	config                 *config.Config
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
		mailService:            mail.NewMailService(nil),
		config:                 config.GetConfig(),
	}
}

// FieldsToDuration convert days/hours/minutes in total minutes
func FieldsToDuration(days, hours, minutes int) int {
	return days*24*60 + hours*60 + minutes
}

func (s *EventService) Create(data *EventCreateDto, user *guard.Claims) (EventCreateResponseDto, error) {
	// Data validation
	data.Name = strings.TrimSpace(data.Name)
	if len(data.Name) < 5 {
		return EventCreateResponseDto{}, errors.New("event name must be at least 5 characters long")
	}
	if data.Description != nil {
		*data.Description = strings.TrimSpace(*data.Description)
		if len(*data.Description) == 0 {
			data.Description = nil
		}
	}

	// Prevent creating events with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return EventCreateResponseDto{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent creating events in the past
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if data.StartsAt.Before(now) {
		return EventCreateResponseDto{}, constants.ERR_EVENT_START_BEFORE_TODAY.Err
	}

	// Prevent creating events with duration less than 1 day
	oneDayAfterStart := data.StartsAt.Add(24 * time.Hour)
	if data.EndsAt.Before(oneDayAfterStart) {
		return EventCreateResponseDto{}, constants.ERR_EVENT_DURATION_TOO_SHORT.Err
	}

	// Validate duration (15 min minimum, 30240 min = 3 weeks maximum)
	duration := FieldsToDuration(data.Days, data.Hours, data.Minutes)
	if duration < 15 || duration > 30240 {
		return EventCreateResponseDto{}, constants.ERR_EVENT_DURATION_TOO_SHORT.Err
	}

	// Create event
	event := model.Event{
		Id:          uuid.New(),
		Name:        data.Name,
		Description: data.Description,
		Duration:    duration,
		StartsAt:    data.StartsAt,
		EndsAt:      data.EndsAt,
		Owner: model.Account{
			Id:       user.Id,
			UserName: user.Username,
		},
		Status: constants.EVENT_STATUS_IN_DECISION,
	}
	if err := s.eventRepository.Create(&event); err != nil {
		return EventCreateResponseDto{}, err
	}

	// Create account_event relation
	accountEvent := model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		_ = s.eventRepository.Delete(event.Id)
		return EventCreateResponseDto{}, err
	}

	// Reload event with full owner data
	if err := s.eventRepository.FindOneById(event.Id, &event); err != nil {
		return EventCreateResponseDto{}, err
	}

	return MapToEventCreateResponseDto(event), nil
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
	if data.Name != nil {
		*data.Name = strings.TrimSpace(*data.Name)
		if len(*data.Name) < 5 {
			return errors.New("event name must be at least 5 characters long")
		}
	}
	if data.Description != nil {
		*data.Description = strings.TrimSpace(*data.Description)
		if len(*data.Description) == 0 {
			data.Description = nil
		}
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
	// If any duration field is provided, recompute the full duration from all three fields
	if data.Days != nil || data.Hours != nil || data.Minutes != nil {
		days, hours, minutes := 0, 0, 0
		if data.Days != nil {
			days = *data.Days
		}
		if data.Hours != nil {
			hours = *data.Hours
		}
		if data.Minutes != nil {
			minutes = *data.Minutes
		}
		duration := FieldsToDuration(days, hours, minutes)
		if duration < 15 || duration > 30240 {
			return constants.ERR_EVENT_DURATION_TOO_SHORT.Err
		}
		event.Duration = duration
		isBreakingSlots = true
	}

	// Update event in repository
	if err := s.eventRepository.Updates(&event); err != nil {
		return err
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
	pagination *lib.Pagination[EventListItemDto],
) error {
	events, total, err := s.eventRepository.FindEventsByAccountId(user.Id, pagination.Limit, pagination.Offset)
	if err != nil {
		return err
	}
	pagination.Total = total

	dtos := make([]EventListItemDto, 0, len(events))
	for i := range events {
		// Update event status if needed
		if _, err := events[i].CheckAndAutoUpdateStatus(s.eventRepository.Updates, nil); err != nil {
			return err
		}
		dtos = append(dtos, MapToEventListItemDto(events[i]))
	}
	pagination.Data = dtos

	return nil
}

func (s *EventService) GetEvent(eventId uuid.UUID, user *guard.Claims) (interface{}, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ERR_EVENT_NOT_FOUND.Err
		}
		return nil, err
	}

	// Update event status if needed
	if _, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, nil); err != nil {
		return nil, err
	}

	// If no user, return event basic info
	if user == nil {
		return MapToEventBasicResponseDto(event), nil
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	notJoined := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !notJoined {
		return nil, err
	}
	if notJoined {
		return MapToEventBasicResponseDto(event), nil
	}

	return MapToEventFullResponseDto(event), nil
}

func (s *EventService) JoinEvent(eventId uuid.UUID, user *guard.Claims) (EventFullResponseDto, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EventFullResponseDto{}, constants.ERR_EVENT_NOT_FOUND.Err
		}
		return EventFullResponseDto{}, err
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	alreadyJoined := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return EventFullResponseDto{}, err
	}
	if alreadyJoined {
		return EventFullResponseDto{}, constants.ERR_EVENT_ALREADY_JOINED.Err
	}

	// Check and update event status if needed
	if hasStatus, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, &[]constants.EventStatus{constants.EVENT_STATUS_IN_DECISION, constants.EVENT_STATUS_UPCOMING}); !hasStatus || err != nil {
		if err != nil {
			return EventFullResponseDto{}, err
		}
		return EventFullResponseDto{}, constants.ERR_EVENT_ENDED.Err
	}

	// Create account_event relation
	accountEvent = model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		return EventFullResponseDto{}, err
	}

	// Reload event with all relations
	if err := s.eventRepository.FindOneById(event.Id, &event); err != nil {
		return EventFullResponseDto{}, err
	}

	return MapToEventFullResponseDto(event), nil
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
