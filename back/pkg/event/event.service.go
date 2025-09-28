package event

import (
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"errors"
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

func (s *EventService) Create(data *EventCreateDto, user *guard.Claims) (model.Event, error) {
	// Prevent creating events with end date before start date
	if data.StartsAt.After(data.EndsAt) {
		return model.Event{}, errors.New("event start date must be before end date")
	}

	// Prevent creating events in the past
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	if data.StartsAt.Before(now) {
		return model.Event{}, errors.New("event start date must be from today")
	}

	// Create event
	event := model.Event{
		Id:       uuid.New(),
		Name:     data.Name,
		Duration: data.Duration,
		StartsAt: data.StartsAt,
		EndsAt:   data.EndsAt,
	}
	if err := s.eventRepository.Create(&event); err != nil {
		return event, err
	}

	// Create account_event relation
	accountEvent := model.AccountEvent{
		AccountId: user.Id,
		EventId:   event.Id,
		IsOwner:   true,
	}
	if err := s.accountEventRepository.Create(&accountEvent); err != nil {
		s.eventRepository.Delete(event.Id)
		return event, err
	}

	event.Accounts = []model.AccountEvent{accountEvent}

	return event, nil
}
