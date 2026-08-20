package slot

import (
	"app/commons/constants"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
	"app/pkg/sse"
	"app/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSlotService(t *testing.T) *SlotService {
	t.Helper()
	testutils.TestDB(t)
	return &SlotService{
		slotRepository:         repository.NewSlotRepository(nil),
		eventRepository:        repository.NewEventRepository(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		accountEventRepository: repository.NewAccountEventRepository(nil),
		sseService:             sse.NewSSEService(),
		mailService:            mail.NewMailService(nil),
	}
}

func createTestAccount(t *testing.T) model.Account {
	t.Helper()
	username := "slot-" + uuid.NewString()
	email := username + "@example.com"
	account := model.Account{Id: uuid.New(), Username: &username, Email: &email}
	require.NoError(t, repository.NewAccountRepository(nil).Create(repository.AccountCreateDto{
		Id: account.Id, Username: &username, Email: &email,
	}, &account))
	return account
}

func createTestEvent(t *testing.T, s *SlotService, ownerId uuid.UUID, participantIds ...uuid.UUID) model.Event {
	t.Helper()
	start := time.Now().UTC().Add(time.Hour).Truncate(time.Minute)
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Slot Test Event",
		Duration: 30,
		StartsAt: start,
		EndsAt:   start.Add(4 * time.Hour),
		OwnerId:  ownerId,
		Status:   constants.EVENT_STATUS_IN_DECISION,
	}
	require.NoError(t, s.eventRepository.Create(&event))

	require.NoError(t, s.accountEventRepository.Create(&model.AccountEvent{AccountId: ownerId, EventId: event.Id}))
	for _, id := range participantIds {
		require.NoError(t, s.accountEventRepository.Create(&model.AccountEvent{AccountId: id, EventId: event.Id}))
	}

	var found model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &found))
	return found
}

func TestNewSlotService_ReusesProvidedInstance(t *testing.T) {
	existing := &SlotService{}
	assert.Same(t, existing, NewSlotService(existing))
}

func TestSlotService_ConfirmSlot_EventEnded(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	event.EndsAt = time.Now().Add(-time.Hour)
	require.NoError(t, s.eventRepository.Updates(&event))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
}

func TestSlotService_ConfirmSlot_AutoFinishUpdateFails(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	event.EndsAt = time.Now().Add(-time.Hour)
	require.NoError(t, s.eventRepository.Updates(&event))

	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
}

func TestSlotService_ConfirmSlot_EventUpdateFails(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.Error(t, err)
}

func TestSlotService_ConfirmSlot_ParticipantsLookupFails_StillSucceeds(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))

	resp, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.NoError(t, err)
	assert.True(t, resp.IsValidated)
}

func TestSlotService_ConfirmSlot_CreateFails(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	testutils.MakeReadOnly(t)

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.Error(t, err)
}

func TestSlotService_ConfirmSlot_NotFound(t *testing.T) {
	s := newTestSlotService(t)
	_, err := s.ConfirmSlot(ConfirmSlotDto{}, uuid.New(), uuid.New())
	assert.ErrorIs(t, err, constants.ERR_SLOT_NOT_FOUND.Err)
}

func TestSlotService_ConfirmSlot_NotOwner(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, uuid.New())
	assert.ErrorIs(t, err, constants.ERR_EVENT_ACCESS_DENIED.Err)
}

func TestSlotService_ConfirmSlot_InvalidStartsAt(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: event.StartsAt.Add(-time.Hour), EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.ErrorIs(t, err, constants.ERR_SLOT_INVALID_STARTS_AT.Err)
}

func TestSlotService_ConfirmSlot_InvalidEndsAt(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt.Add(time.Hour)}, slotEntity.Id, owner.Id)
	assert.ErrorIs(t, err, constants.ERR_SLOT_INVALID_ENDS_AT.Err)
}

func TestSlotService_ConfirmSlot_Success(t *testing.T) {
	s := newTestSlotService(t)
	called := testutils.StubSMTPAwait(t, &s.mailService.SendMailFunc)

	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	resp, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.NoError(t, err)
	assert.True(t, resp.IsValidated)

	var updatedEvent model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &updatedEvent))
	assert.Equal(t, constants.EVENT_STATUS_UPCOMING, updatedEvent.Status)

	// Wait for the async confirmation email goroutine (owner is the sole participant).
	testutils.AwaitSMTP(t, called)
}

func TestSlotService_RemoveValidatedSlot_EventUpdateFails(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	s.eventRepository = repository.NewEventRepository(testutils.ClosedDB(t))

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.Error(t, err)
}

func TestSlotService_RemoveValidatedSlot_ParticipantsLookupFails_StillSucceeds(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	s.accountEventRepository = repository.NewAccountEventRepository(testutils.ClosedDB(t))

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.NoError(t, err)
}

func TestSlotService_RemoveValidatedSlot_DeleteFails(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	testutils.MakeReadOnly(t)

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.Error(t, err)
}

func TestSlotService_RemoveValidatedSlot_NotFound(t *testing.T) {
	s := newTestSlotService(t)
	err := s.RemoveValidatedSlot(uuid.New(), uuid.New())
	assert.ErrorIs(t, err, constants.ERR_SLOT_NOT_FOUND.Err)
}

func TestSlotService_RemoveValidatedSlot_NotOwner(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	err := s.RemoveValidatedSlot(slotEntity.Id, uuid.New())
	assert.ErrorIs(t, err, constants.ERR_EVENT_ACCESS_DENIED.Err)
}

func TestSlotService_RemoveValidatedSlot_NotValidated(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.ErrorIs(t, err, constants.ERR_SLOT_NOT_FOUND.Err)
}

func TestSlotService_RemoveValidatedSlot_Success(t *testing.T) {
	s := newTestSlotService(t)
	called := testutils.StubSMTPAwait(t, &s.mailService.SendMailFunc)

	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.NoError(t, err)
	// RemoveValidatedSlot also spawns `go s.LoadSlots(...)`, which reuses this
	// test's dedicated connection; poll (in this goroutine) instead of
	// blindly sleeping before querying again.
	var updatedEvent model.Event
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.eventRepository.FindOneById(event.Id, &updatedEvent); err != nil {
			return false
		}
		return updatedEvent.Status == constants.EVENT_STATUS_IN_DECISION
	})
	assert.Equal(t, constants.EVENT_STATUS_IN_DECISION, updatedEvent.Status)

	// Wait for the async cancellation email goroutine (owner is the sole participant).
	testutils.AwaitSMTP(t, called)
}

func TestSlotService_LoadSlots_NotEnoughParticipants(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	assert.Len(t, slots, 0)
}

func TestSlotService_LoadSlots_CreatesCommonSlots(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)

	availRepo := repository.NewAvailabilityRepository(nil)
	require.NoError(t, availRepo.Create(&model.Availability{
		Id: uuid.New(), AccountId: owner.Id, EventId: event.Id,
		StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour),
	}))
	require.NoError(t, availRepo.Create(&model.Availability{
		Id: uuid.New(), AccountId: participant.Id, EventId: event.Id,
		StartsAt: event.StartsAt.Add(15 * time.Minute), EndsAt: event.StartsAt.Add(75 * time.Minute),
	}))

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	assert.Greater(t, len(slots), 0)
}

func TestSlotService_LoadSlots_EventNotFound_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	// Should not panic when the event doesn't exist.
	s.LoadSlots(uuid.New())
}

func TestSlotService_LoadSlots_EventLocked_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	event.EndsAt = time.Now().Add(-time.Hour)
	require.NoError(t, s.eventRepository.Updates(&event))

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	assert.Len(t, slots, 0)
}

func TestSlotService_LoadSlots_DeleteSlotsFails_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)

	s.slotRepository = repository.NewSlotRepository(testutils.ClosedDB(t))

	// Should not panic; the deletion error is logged and swallowed.
	s.LoadSlots(event.Id)
}

func TestSlotService_LoadSlots_AvailabilitiesLookupFails_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)

	s.availabilityRepository = repository.NewAvailabilityRepository(testutils.ClosedDB(t))

	// Should not panic; the lookup error is logged and swallowed.
	s.LoadSlots(event.Id)
}

func TestSlotService_LoadSlots_NoCommonSlots_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)

	availRepo := repository.NewAvailabilityRepository(nil)
	require.NoError(t, availRepo.Create(&model.Availability{
		Id: uuid.New(), AccountId: owner.Id, EventId: event.Id,
		StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute),
	}))
	require.NoError(t, availRepo.Create(&model.Availability{
		Id: uuid.New(), AccountId: participant.Id, EventId: event.Id,
		StartsAt: event.StartsAt.Add(2 * time.Hour), EndsAt: event.StartsAt.Add(3 * time.Hour),
	}))

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	assert.Len(t, slots, 0)
}

func TestFindIntersectingTimeSlots_SingleUser_ReturnsEmpty(t *testing.T) {
	s := newTestSlotService(t)
	result := s.findIntersectingTimeSlots(map[uuid.UUID][]TimeSlot{
		uuid.New(): {{StartsAt: time.Now(), EndsAt: time.Now().Add(time.Hour)}},
	}, 15*time.Minute)
	assert.Empty(t, result)
}
