package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/helpers"
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
	"github.com/rs/zerolog/log"
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
		eventRepository:        repository.NewEventRepository(nil),
		accountEventRepository: repository.NewAccountEventRepository(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		slotRepository:         repository.NewSlotRepository(nil),
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
	var fieldErrors helpers.FieldErrors
	data.Name = strings.TrimSpace(data.Name)
	if len(data.Name) < 5 {
		fieldErrors.Add(constants.ERR_FIELD_TOO_SHORT, "name")
	}
	if data.Description != nil {
		*data.Description = strings.TrimSpace(*data.Description)
		if len(*data.Description) == 0 {
			data.Description = nil
		}
	}

	validateEventDates(&fieldErrors, data.StartsAt, data.EndsAt, true)

	duration := FieldsToDuration(data.Days, data.Hours, data.Minutes)
	validateEventDuration(&fieldErrors, duration)

	if err := fieldErrors.Err(); err != nil {
		return EventCreateResponseDto{}, err
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
			Username: user.Username,
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

// validateEventDates collects the date rules shared by event creation and update
func validateEventDates(fieldErrors *helpers.FieldErrors, startsAt, endsAt time.Time, checkStartsAt bool) {
	// Prevent end date before start date
	if startsAt.After(endsAt) {
		fieldErrors.Add(constants.ERR_EVENT_START_AFTER_END, "endsAt")
	}
	// Prevent duration less than 1 day
	if endsAt.Before(startsAt.Add(24 * time.Hour)) {
		fieldErrors.Add(constants.ERR_EVENT_DURATION_TOO_SHORT, "endsAt")
	}
	// Prevent dates before today
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if checkStartsAt && startsAt.Before(today) {
		fieldErrors.Add(constants.ERR_EVENT_START_BEFORE_TODAY, "startsAt")
	}
	if endsAt.Before(today) {
		fieldErrors.Add(constants.ERR_EVENT_START_BEFORE_TODAY, "endsAt")
	}
}

// validateEventDuration checks the duration (15 min minimum, 30240 min = 3 weeks maximum)
func validateEventDuration(fieldErrors *helpers.FieldErrors, duration int) {
	if duration < 15 || duration > 30240 {
		fieldErrors.Add(constants.ERR_EVENT_DURATION_TOO_SHORT, "duration")
	}
}

// resolveEventDates merges the DTO dates with the current event dates
func resolveEventDates(event *model.Event, startsAtDto, endsAtDto *time.Time) (startsAt, endsAt time.Time) {
	startsAt, endsAt = event.StartsAt, event.EndsAt
	if startsAtDto != nil {
		startsAt = *startsAtDto
	}
	if endsAtDto != nil {
		endsAt = *endsAtDto
	}
	return startsAt, endsAt
}

// applyEventDates sets already validated dates, unless they conflict with the validated slot
func applyEventDates(event *model.Event, startsAt, endsAt time.Time) error {
	if validatedSlot := event.GetValidatedSlot(); validatedSlot != nil && (endsAt.Before(validatedSlot.EndsAt) || startsAt.After(validatedSlot.StartsAt)) {
		// Prevent updating events to end date before already validated slot
		return constants.ERR_VALIDATED_SLOT_CANNOT_BE_MODIFIED.Err
	}

	event.StartsAt = startsAt
	event.EndsAt = endsAt

	return nil
}

// SetEventDatesFromDto validates and sets the event dates from the provided DTO values.
func SetEventDatesFromDto(event *model.Event, startsAtDto, endsAtDto *time.Time) error {
	if event == nil {
		return errors.New("event pointer is nil")
	}
	if startsAtDto == nil && endsAtDto == nil {
		return nil
	}

	startsAt, endsAt := resolveEventDates(event, startsAtDto, endsAtDto)

	var fieldErrors helpers.FieldErrors
	validateEventDates(&fieldErrors, startsAt, endsAt, startsAtDto != nil)
	if err := fieldErrors.Err(); err != nil {
		return err
	}

	return applyEventDates(event, startsAt, endsAt)
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
	var fieldErrors helpers.FieldErrors
	if data.Name != nil {
		*data.Name = strings.TrimSpace(*data.Name)
		if len(*data.Name) < 5 {
			fieldErrors.Add(constants.ERR_FIELD_TOO_SHORT, "name")
		}
	}
	if data.Description != nil {
		*data.Description = strings.TrimSpace(*data.Description)
		if len(*data.Description) == 0 {
			data.Description = nil
		}
	}

	datesChanged := data.StartsAt != nil || data.EndsAt != nil
	startsAt, endsAt := resolveEventDates(&event, data.StartsAt, data.EndsAt)
	if datesChanged {
		validateEventDates(&fieldErrors, startsAt, endsAt, data.StartsAt != nil)
	}

	// If any duration field is provided, recompute from current duration to preserve unset fields
	durationChanged := data.Days != nil || data.Hours != nil || data.Minutes != nil
	duration := event.Duration
	if durationChanged {
		current := durationToFields(event.Duration)
		days, hours, minutes := current.Days, current.Hours, current.Minutes
		if data.Days != nil {
			days = *data.Days
		}
		if data.Hours != nil {
			hours = *data.Hours
		}
		if data.Minutes != nil {
			minutes = *data.Minutes
		}
		duration = FieldsToDuration(days, hours, minutes)
		validateEventDuration(&fieldErrors, duration)
	}

	if err := fieldErrors.Err(); err != nil {
		return err
	}

	// Update fields if provided
	if datesChanged {
		if err := applyEventDates(&event, startsAt, endsAt); err != nil {
			return err
		}
	}
	if data.Name != nil {
		event.Name = *data.Name
	}
	if data.Description != nil {
		event.Description = data.Description
	}
	event.Duration = duration
	isBreakingSlots := datesChanged || durationChanged

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
	s.slotService.LoadSlotsAsync(eventId)

	return nil
}

func (s *EventService) Cancel(eventId uuid.UUID, user *guard.Claims) error {
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

	// Prevent canceling closed events
	if hasStatus, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, &[]constants.EventStatus{constants.EVENT_STATUS_IN_DECISION, constants.EVENT_STATUS_UPCOMING}); !hasStatus || err != nil {
		if err != nil {
			return err
		}
		return constants.ERR_EVENT_CLOSED_CANNOT_BE_CANCELED.Err
	}

	// Update event status to canceled
	event.Status = constants.EVENT_STATUS_CANCELED
	if err := s.eventRepository.Updates(&event); err != nil {
		return err
	}

	// Get participants to notify
	var participants []model.Account
	if err := s.accountEventRepository.FindAccountsByEventId(event.Id, &participants); err != nil {
		log.Error().Err(err).Str("eventId", event.Id.String()).Msg("Failed to get participants for event cancellation mail")
	}

	// Send cancellation emails to all participants (including owner)
	for _, participant := range participants {
		s.mailService.SendEventCanceledEmail(participant, event)
	}

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

// GetEventSummary returns basic event info, accessible without authentication
func (s *EventService) GetEventSummary(eventId uuid.UUID) (EventBasicResponseDto, error) {
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EventBasicResponseDto{}, constants.ERR_EVENT_NOT_FOUND.Err
		}
		return EventBasicResponseDto{}, err
	}

	if _, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, nil); err != nil {
		return EventBasicResponseDto{}, err
	}

	return MapToEventBasicResponseDto(event), nil
}

// GetEvent returns full event data, restricted to authenticated members
func (s *EventService) GetEvent(eventId uuid.UUID, user *guard.Claims) (EventFullResponseDto, error) {
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EventFullResponseDto{}, constants.ERR_EVENT_NOT_FOUND.Err
		}
		return EventFullResponseDto{}, err
	}

	if _, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, nil); err != nil {
		return EventFullResponseDto{}, err
	}

	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return EventFullResponseDto{}, constants.ERR_EVENT_NOT_FOUND.Err
	}
	if err != nil {
		return EventFullResponseDto{}, err
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
		return helpers.NewFieldError(constants.ERR_INVALID_COLOR_FORMAT, "color")
	}

	// Update color property
	accountEvent.Color = &data.Color
	if err := s.accountEventRepository.Updates(&accountEvent); err != nil {
		return err
	}

	return nil
}
