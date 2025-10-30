package event

import (
	"app/commons/constants"
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"time"

	"github.com/google/uuid"
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
		s.eventRepository.Delete(event.Id)
		return EventResponse{}, err
	}
	if err := s.accountEventRepository.FindByAccountAndEventId(user.Id, event.Id, &accountEvent); err != nil {
		s.eventRepository.Delete(event.Id)
		return EventResponse{}, err
	}

	return EventResponse{
		Event: event,
		Accounts: []model.Account{{
			Id:       accountEvent.Account.Id,
			UserName: accountEvent.Account.UserName,
		}},
	}, nil
}

func (s *EventService) GetMyEvents(user *guard.Claims) ([]EventResponse, error) {
	events := []EventResponse{}

	// Find all account_events for this user
	var accountEvents []model.AccountEvent
	if err := s.accountEventRepository.FindByAccountId(user.Id, &accountEvents); err != nil {
		return []EventResponse{}, err
	}

	eventIds := make([]uuid.UUID, 0, len(accountEvents))
	for _, accountEvent := range accountEvents {
		eventIds = append(eventIds, accountEvent.EventId)
	}

	// Find all account_events related to these events
	accountEvents = []model.AccountEvent{}
	if err := s.accountEventRepository.FindByIds(eventIds, &accountEvents); err != nil {
		return []EventResponse{}, err
	}

	// Group accounts by event
	for _, accountEvent := range accountEvents {
		accounts := []model.Account{}
		// Find all accounts for this event
		for _, ae := range accountEvents {
			if ae.EventId == accountEvent.EventId {
				accounts = append(accounts, model.Account{
					Id:       ae.Account.Id,
					UserName: ae.Account.UserName,
				})
			}
		}

		accountEvent.Event.Owner = model.Account{
			Id:       accountEvent.Event.Owner.Id,
			UserName: accountEvent.Event.Owner.UserName,
		}

		events = append(events, EventResponse{
			Event:    accountEvent.Event,
			Accounts: accounts,
		})
	}

	return events, nil
}
