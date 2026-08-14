package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"app/pkg/slot"
	"app/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEventService(t *testing.T) *EventService {
	t.Helper()
	testutils.TestDB(t)
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

// createEndedEventForOwner creates an event whose EndsAt is already in the
// past (and status not yet FINISHED), so CheckAndAutoUpdateStatus attempts
// to auto-transition it via eventRepository.Updates.
func createEndedEventForOwner(t *testing.T, s *EventService, ownerId uuid.UUID) model.Event {
	t.Helper()
	start := time.Now().UTC().Add(-4 * time.Hour)
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Ended Event",
		Duration: 60,
		StartsAt: start,
		EndsAt:   start.Add(time.Hour),
		OwnerId:  ownerId,
		Status:   constants.EVENT_STATUS_IN_DECISION,
	}
	require.NoError(t, s.eventRepository.Create(&event))
	require.NoError(t, s.accountEventRepository.Create(&model.AccountEvent{AccountId: ownerId, EventId: event.Id}))

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	return found
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

func TestNewEventService_Nil_BuildsRealDependencies(t *testing.T) {
	s := NewEventService(nil)
	assert.NotNil(t, s.eventRepository)
	assert.NotNil(t, s.accountEventRepository)
	assert.NotNil(t, s.availabilityRepository)
	assert.NotNil(t, s.slotRepository)
	assert.NotNil(t, s.slotService)
	assert.NotNil(t, s.signinService)
	assert.NotNil(t, s.mailService)
}

func TestEventService_Create_NameTooShort(t *testing.T) {
	s := newTestEventService(t)
	dto := validEventCreateDto()
	dto.Name = "ab"

	_, err := s.Create(dto, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_Create_BlankDescriptionIsCleared(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	username := "owner"
	dto := validEventCreateDto()
	blank := "   "
	dto.Description = &blank

	resp, err := s.Create(dto, &guard.Claims{Id: owner, Username: &username})
	require.NoError(t, err)

	var event model.Event
	require.NoError(t, s.eventRepository.FindOneById(resp.Id, &event))
	assert.Nil(t, event.Description)
}

func TestEventService_Create_InvalidDuration(t *testing.T) {
	s := newTestEventService(t)
	dto := validEventCreateDto()
	// Date range is valid (>= 1 day) but the explicit duration fields sum to 0.
	dto.Days, dto.Hours, dto.Minutes = 0, 0, 0

	_, err := s.Create(dto, &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_EVENT_DURATION_TOO_SHORT.Err)
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

func TestEventService_Create_RepositoryCreateError(t *testing.T) {
	s := newTestEventService(t)
	dto := validEventCreateDto()
	testutils.MakeReadOnly(t)

	_, err := s.Create(dto, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_Create_AccountEventCreateError_RollsBackEvent(t *testing.T) {
	s := newTestEventService(t)
	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))
	dto := validEventCreateDto()
	username := "owner"

	_, err := s.Create(dto, &guard.Claims{Id: uuid.New(), Username: &username})
	assert.Error(t, err)
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

func TestEventService_Update_EventRepositoryError(t *testing.T) {
	s := newTestEventService(t)
	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	err := s.Update(uuid.New(), &EventUpdateDto{}, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
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
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.eventRepository.FindOneById(event.Id, &found); err != nil {
			return false
		}
		return found.StartsAt.Equal(newStart)
	})
	assert.True(t, found.StartsAt.Equal(newStart))
}

func TestEventService_Update_DescriptionCleared(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	blank := "   "
	err := s.Update(event.Id, &EventUpdateDto{Description: &blank}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	assert.Nil(t, found.Description)
}

func TestEventService_Update_DescriptionSet(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	description := "A real description"
	err := s.Update(event.Id, &EventUpdateDto{Description: &description}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	require.NotNil(t, found.Description)
	assert.Equal(t, description, *found.Description)
}

func TestEventService_Update_RepositoryUpdatesError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	name := "Renamed"
	err := s.Update(event.Id, &EventUpdateDto{Name: &name}, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_Update_SlotRepositoryDeleteError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	s.slotRepository = repository.NewSlotRepository(testutils.ClosedDB(t))

	newStart := event.StartsAt.Add(time.Hour)
	newEnd := event.EndsAt.Add(time.Hour)
	err := s.Update(event.Id, &EventUpdateDto{StartsAt: &newStart, EndsAt: &newEnd}, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_Update_AvailabilityRepositoryDeleteError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	s.availabilityRepository = repository.NewAvailabilityRepository(testutils.ClosedDB(t))

	newStart := event.StartsAt.Add(time.Hour)
	newEnd := event.EndsAt.Add(time.Hour)
	err := s.Update(event.Id, &EventUpdateDto{StartsAt: &newStart, EndsAt: &newEnd}, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_Update_InvalidDates(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	// End before start -> SetEventDatesFromDto rejects it.
	newStart := event.EndsAt.Add(time.Hour)
	err := s.Update(event.Id, &EventUpdateDto{StartsAt: &newStart}, &guard.Claims{Id: owner})
	assert.ErrorIs(t, err, constants.ERR_EVENT_START_AFTER_END.Err)
}

func TestEventService_Update_DaysOnly_RecalculatesSlots(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	days := 2
	err := s.Update(event.Id, &EventUpdateDto{Days: &days}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.eventRepository.FindOneById(event.Id, &found); err != nil {
			return false
		}
		return found.Duration != event.Duration
	})
	assert.NotEqual(t, event.Duration, found.Duration)
}

func TestEventService_Update_MinutesOnly_RecalculatesSlots(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	minutes := 30
	err := s.Update(event.Id, &EventUpdateDto{Minutes: &minutes}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.eventRepository.FindOneById(event.Id, &found); err != nil {
			return false
		}
		return found.Duration != event.Duration
	})
	assert.NotEqual(t, event.Duration, found.Duration)
}

func TestEventService_Update_InvalidDuration(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	days, hours, minutes := 0, 0, 0
	err := s.Update(event.Id, &EventUpdateDto{Days: &days, Hours: &hours, Minutes: &minutes}, &guard.Claims{Id: owner})
	assert.ErrorIs(t, err, constants.ERR_EVENT_DURATION_TOO_SHORT.Err)
}

func TestEventService_Update_DurationOnly_RecalculatesSlots(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)

	hours := 2
	err := s.Update(event.Id, &EventUpdateDto{Hours: &hours}, &guard.Claims{Id: owner})
	assert.NoError(t, err)

	var found model.Event
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.eventRepository.FindOneById(event.Id, &found); err != nil {
			return false
		}
		return found.Duration != event.Duration
	})
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

func TestEventService_GetUserEvents_RepositoryError(t *testing.T) {
	s := newTestEventService(t)
	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	pagination := &lib.Pagination[EventListItemDto]{Limit: 10}
	err := s.GetUserEvents(&guard.Claims{Id: uuid.New()}, pagination)
	assert.Error(t, err)
}

func TestEventService_GetUserEvents_AutoFinishUpdateError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	createEndedEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	pagination := &lib.Pagination[EventListItemDto]{Limit: 10}
	err := s.GetUserEvents(&guard.Claims{Id: owner}, pagination)
	assert.Error(t, err)
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

func TestEventService_GetEventSummary_RepositoryError(t *testing.T) {
	s := newTestEventService(t)
	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	_, err := s.GetEventSummary(uuid.New())
	assert.Error(t, err)
}

func TestEventService_GetEventSummary_AutoFinishUpdateError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createEndedEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	_, err := s.GetEventSummary(event.Id)
	assert.Error(t, err)
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

func TestEventService_GetEvent_RepositoryError(t *testing.T) {
	s := newTestEventService(t)
	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	_, err := s.GetEvent(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_GetEvent_AutoFinishUpdateError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createEndedEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	_, err := s.GetEvent(event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_GetEvent_AccountEventRepositoryError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))

	_, err := s.GetEvent(event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
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

	newMember := createTestOwner(t)
	dto, err := s.JoinEvent(event.Id, &guard.Claims{Id: newMember})
	assert.NoError(t, err)
	assert.Equal(t, event.Name, dto.Name)

	var accountEvent model.AccountEvent
	require.NoError(t, s.accountEventRepository.FindByAccountAndEventId(newMember, event.Id, &accountEvent))
}

func TestEventService_JoinEvent_EventRepositoryError(t *testing.T) {
	s := newTestEventService(t)
	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	_, err := s.JoinEvent(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_JoinEvent_AccountEventLookupError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))

	_, err := s.JoinEvent(event.Id, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_JoinEvent_AutoFinishUpdateError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createEndedEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	_, err := s.JoinEvent(event.Id, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestEventService_JoinEvent_AccountEventCreateError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	_, err := s.JoinEvent(event.Id, &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
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

func TestEventService_UpdateProfile_LookupError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))

	err := s.UpdateProfile(&EventProfileDto{Color: "#AABBCC"}, event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestEventService_UpdateProfile_RepositoryUpdatesError(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	testutils.MakeReadOnly(t)

	err := s.UpdateProfile(&EventProfileDto{Color: "#AABBCC"}, event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestNewEventService_ReusesProvidedInstance(t *testing.T) {
	existing := &EventService{}
	assert.Same(t, existing, NewEventService(existing))
}
