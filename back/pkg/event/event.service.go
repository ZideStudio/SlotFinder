package event

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService struct {
	eventRepository        *repository.EventRepository
	accountEventRepository *repository.AccountEventRepository
	signinService          *signin.SigninService
}

func NewEventService(service *EventService) *EventService {
	if service != nil {
		return service
	}

	return &EventService{
		eventRepository:        &repository.EventRepository{},
		accountEventRepository: &repository.AccountEventRepository{},
		signinService:          signin.NewSigninService(nil),
	}
}

func (s *EventService) Create(data *EventCreateDto, user *guard.Claims) (EventResponse, error) {
	// Prevent creating events with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return EventResponse{}, constants.ERR_EVENT_START_AFTER_END.Err
	}

	// Prevent creating events in the past
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if data.StartsAt.Before(now) {
		return EventResponse{}, constants.ERR_EVENT_START_BEFORE_TODAY.Err
	}

	// Prevent creating events with duration less than 1 day
	oneDayAfterStart := data.StartsAt.Add(24 * time.Hour)
	if data.EndsAt.Before(oneDayAfterStart) {
		return EventResponse{}, constants.ERR_EVENT_DURATION_TOO_SHORT.Err
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
	}
	if err := s.eventRepository.Create(&event); err != nil {
		return EventResponse{}, err
	}

	// Create account_event relation
	accountEvent := model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		_ = s.eventRepository.Delete(event.Id)
		return EventResponse{}, err
	}
	if err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent); err != nil {
		_ = s.eventRepository.Delete(event.Id)
		return EventResponse{}, err
	}

	return EventResponse{
		Event:    event,
		Accounts: []model.Account{accountEvent.Account.Sanitized()},
	}, nil
}

func (s *EventService) getEventResponseFromEvents(eventIds []uuid.UUID) ([]EventResponse, error) {
	// Find all account_events for these events to get all accounts
	var allAccountEvents []model.AccountEvent
	if err := s.accountEventRepository.FindByIds(eventIds, &allAccountEvents); err != nil {
		return nil, err
	}

	// Group accounts by event ID
	type eventGroup struct {
		event    model.Event
		accounts []model.Account
	}
	eventMap := make(map[uuid.UUID]*eventGroup)
	for _, ae := range allAccountEvents {
		eg, ok := eventMap[ae.EventId]
		if !ok {
			eg = &eventGroup{
				event: ae.Event,
			}
			eventMap[ae.EventId] = eg
		}

		eg.accounts = append(eg.accounts, ae.Account.Sanitized())
	}

	// Build final event response
	events := make([]EventResponse, 0, len(eventMap))
	for _, eg := range eventMap {
		eg.event.Owner = eg.event.Owner.Sanitized()

		events = append(events, EventResponse{
			Event:    eg.event,
			Accounts: eg.accounts,
		})
	}

	return events, nil
}

func (s *EventService) GetUserEvents(user *guard.Claims) ([]EventResponse, error) {
	events := []EventResponse{}

	// Find all account_events for this user
	var myAccountEvents []model.AccountEvent
	if err := s.accountEventRepository.FindByAccountId(user.Id, &myAccountEvents); err != nil {
		return events, err
	}

	// Collect event IDs the user is associated with
	eventIds := make([]uuid.UUID, 0, len(myAccountEvents))
	for _, accountEvent := range myAccountEvents {
		eventIds = append(eventIds, accountEvent.EventId)
	}
	if len(eventIds) == 0 {
		return events, nil
	}

	return s.getEventResponseFromEvents(eventIds)
}

func (s *EventService) GetEvent(eventId uuid.UUID, user *guard.Claims) (EventResponse, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EventResponse{}, constants.ERR_EVENT_NOT_FOUND.Err
		}

		return EventResponse{}, err
	}
	event.Owner = event.Owner.Sanitized()

	// If no user, return event basic info
	if user == nil {
		return EventResponse{
			Event:    event,
			Accounts: nil,
		}, nil
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	notJoined := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !notJoined {
		return EventResponse{}, err
	}

	// Return event info
	eventResponse, err := s.getEventResponseFromEvents([]uuid.UUID{event.Id})
	if err != nil {
		return EventResponse{}, err
	}
	if len(eventResponse) == 0 {
		return EventResponse{}, errors.New("failed to get event response")
	}

	return eventResponse[0], nil
}

func (s *EventService) JoinEvent(eventId uuid.UUID, user *guard.Claims) (EventResponse, error) {
	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EventResponse{}, constants.ERR_EVENT_NOT_FOUND.Err
		}

		return EventResponse{}, err
	}

	// Check if user already joined the event
	var accountEvent model.AccountEvent
	err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent)
	alreadyJoin := !errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && alreadyJoin {
		return EventResponse{}, err
	}
	if alreadyJoin {
		return EventResponse{}, constants.ERR_EVENT_ALREADY_JOINED.Err
	}

	// Create account_event relation
	accountEvent = model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		return EventResponse{}, err
	}
	if err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent); err != nil {
		return EventResponse{}, err
	}

	eventResponse, err := s.getEventResponseFromEvents([]uuid.UUID{event.Id})
	if err != nil {
		return EventResponse{}, err
	}
	if len(eventResponse) == 0 {
		return EventResponse{}, errors.New("failed to get event response")
	}

	return eventResponse[0], nil
}
