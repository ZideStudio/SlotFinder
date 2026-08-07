package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"app/pkg/slot"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEventService(t *testing.T) *EventService {
	t.Helper()
	return &EventService{
		eventRepository:        repository.NewEventRepository(nil),
		accountEventRepository: repository.NewAccountEventRepository(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		slotRepository:         repository.NewSlotRepository(nil),
		slotService:            slot.NewSlotService(nil),
		signinService:          signin.NewSigninService(nil),
	}
}

func createTestOwner(t *testing.T) uuid.UUID {
	t.Helper()
	id := uuid.New()
	username := "event-owner-" + uuid.NewString()
	require.NoError(t, repository.NewAccountRepository(nil).Create(repository.AccountCreateDto{
		Id: id, UserName: &username,
	}, &model.Account{}))
	return id
}

func validEventCreateDto() *EventCreateDto {
	start := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Hour)
	return &EventCreateDto{
		Name:     "My Test Event",
		Days:     1,
		StartsAt: start,
		EndsAt:   start.Add(48 * time.Hour),
	}
}

func TestEventService_Create_Success(t *testing.T) {
	s := newTestEventService(t)
	ownerId := createTestOwner(t)
	username := "owner"
	dto := validEventCreateDto()

	resp, err := s.Create(dto, &guard.Claims{Id: ownerId, Username: &username})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, resp.Id)

	var accountEvent model.AccountEvent
	require.NoError(t, s.accountEventRepository.FindByAccountAndEventId(ownerId, resp.Id, &accountEvent))
}

func TestEventService_Update_NotFound(t *testing.T) {
	s := newTestEventService(t)
	err := s.Update(uuid.New(), &EventUpdateDto{}, &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func createTestEventForOwner(t *testing.T, s *EventService, ownerId uuid.UUID) model.Event {
	t.Helper()
	username := "owner"
	resp, err := s.Create(validEventCreateDto(), &guard.Claims{Id: ownerId, Username: &username})
	require.NoError(t, err)

	var event model.Event
	require.NoError(t, s.eventRepository.FindOneById(resp.Id, &event))
	return event
}

func TestEventService_Update_AccessDenied(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	err := s.Update(event.Id, &EventUpdateDto{}, &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_ACCESS_DENIED.Err)
}

func TestEventService_Update_NameTooShort(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	name := "ab"
	err := s.Update(event.Id, &EventUpdateDto{Name: &name}, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_Update_NonBreakingChange(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	name := "Updated Event Name"
	err := s.Update(event.Id, &EventUpdateDto{Name: &name}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	assert.Equal(t, name, found.Name)
}

func TestEventService_Update_BreakingChange_RecalculatesSlots(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	newStart := event.StartsAt.Add(time.Hour)
	newEnd := event.EndsAt.Add(time.Hour)
	err := s.Update(event.Id, &EventUpdateDto{StartsAt: &newStart, EndsAt: &newEnd}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	assert.True(t, found.StartsAt.Equal(newStart))
}

func TestEventService_Update_DurationOnly_RecalculatesSlots(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	hours := 2
	err := s.Update(event.Id, &EventUpdateDto{Hours: &hours}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	assert.NotEqual(t, event.Duration, found.Duration)
}

func TestEventService_GetUserEvents(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	createTestEventForOwner(t, s, owner)
	createTestEventForOwner(t, s, owner)

	pagination := &lib.Pagination[EventListItemDto]{Limit: 10}
	err := s.GetUserEvents(&guard.Claims{Id: owner}, pagination)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), pagination.Total)
	assert.Len(t, pagination.Data, 2)
}

func TestEventService_GetEventSummary_NotFound(t *testing.T) {
	s := newTestEventService(t)
	_, err := s.GetEventSummary(uuid.New())
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestEventService_GetEventSummary_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	dto, err := s.GetEventSummary(event.Id)
	assert.NoError(t, err)
	assert.Equal(t, event.Name, dto.Name)
}

func TestEventService_GetEvent_NotFound(t *testing.T) {
	s := newTestEventService(t)
	_, err := s.GetEvent(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestEventService_GetEvent_NotAMember(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	_, err := s.GetEvent(event.Id, &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestEventService_GetEvent_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	dto, err := s.GetEvent(event.Id, &guard.Claims{Id: owner})
	assert.NoError(t, err)
	assert.Equal(t, event.Name, dto.Name)
}

func TestEventService_JoinEvent_NotFound(t *testing.T) {
	s := newTestEventService(t)
	_, err := s.JoinEvent(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestEventService_JoinEvent_AlreadyJoined(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	_, err := s.JoinEvent(event.Id, &guard.Claims{Id: owner})
	assert.ErrorIs(t, err, constants.ERR_EVENT_ALREADY_JOINED.Err)
}

func TestEventService_JoinEvent_EventEnded(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	event.EndsAt = time.Now().Add(-time.Hour)
	require.NoError(t, s.eventRepository.Updates(&event))

	newMember := uuid.New()
	_, err := s.JoinEvent(event.Id, &guard.Claims{Id: newMember})
	assert.ErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
}

func TestEventService_JoinEvent_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	newMember := uuid.New()
	dto, err := s.JoinEvent(event.Id, &guard.Claims{Id: newMember})
	assert.NoError(t, err)
	assert.Equal(t, event.Name, dto.Name)

	var accountEvent model.AccountEvent
	require.NoError(t, s.accountEventRepository.FindByAccountAndEventId(newMember, event.Id, &accountEvent))
}

func TestEventService_UpdateProfile_NotAMember(t *testing.T) {
	s := newTestEventService(t)
	err := s.UpdateProfile(&EventProfileDto{Color: "#FFFFFF"}, uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestEventService_UpdateProfile_InvalidColor(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	err := s.UpdateProfile(&EventProfileDto{Color: "not-a-color"}, event.Id, &guard.Claims{Id: owner})
	assert.ErrorIs(t, err, constants.ERR_INVALID_COLOR_FORMAT.Err)
}

func TestEventService_UpdateProfile_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	err := s.UpdateProfile(&EventProfileDto{Color: "#AABBCC"}, event.Id, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var accountEvent model.AccountEvent
	require.NoError(t, s.accountEventRepository.FindByAccountAndEventId(owner, event.Id, &accountEvent))
	require.NotNil(t, accountEvent.Color)
	assert.Equal(t, "#AABBCC", *accountEvent.Color)
}

func TestNewEventService_ReusesProvidedInstance(t *testing.T) {
	existing := &EventService{}
	assert.Same(t, existing, NewEventService(existing))
}
