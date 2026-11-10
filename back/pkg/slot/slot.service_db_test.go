package slot

import (
	"app/commons/constants"
	"app/commons/helpers"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
	"app/pkg/sse"
	"app/testutils"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSlotService(t *testing.T) *SlotService {
	t.Helper()
	testutils.TestDB(t)
	s := &SlotService{
		slotRepository:         repository.NewSlotRepository(nil),
		eventRepository:        repository.NewEventRepository(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		accountEventRepository: repository.NewAccountEventRepository(nil),
		sseService:             sse.NewSSEService(),
		mailService:            mail.NewMailService(nil),
	}
	// Wait for trailing LoadSlotsAsync goroutines before rollback (t.Cleanup
	// is LIFO, so this runs first).
	t.Cleanup(s.WaitAsyncLoads)
	return s
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

func createTestAvailability(t *testing.T, accountId uuid.UUID, eventId uuid.UUID, startsAt time.Time, endsAt time.Time) model.Availability {
	t.Helper()
	availability := model.Availability{Id: uuid.New(), AccountId: accountId, EventId: eventId, StartsAt: startsAt, EndsAt: endsAt}
	require.NoError(t, repository.NewAvailabilityRepository(nil).Create(&availability))
	return availability
}

// readSlotUpdates returns the next SSE message sent to the client
func readSlotUpdates(t *testing.T, client *sse.SSEClient) []sse.SSESlotUpdateMessage {
	t.Helper()
	select {
	case msg := <-client.Channel:
		var updates []sse.SSESlotUpdateMessage
		require.NoError(t, json.Unmarshal(msg, &updates))
		return updates
	case <-time.After(time.Second):
		t.Fatal("expected an SSE message")
		return nil
	}
}

// assertSlotUpdate checks that an SSE message describes the slot with the given action
func assertSlotUpdate(t *testing.T, action constants.SlotAction, slot model.Slot, update sse.SSESlotUpdateMessage) {
	t.Helper()
	assert.Equal(t, action, update.Action)
	assert.Equal(t, slot.Id, update.Id)
	assert.True(t, slot.StartsAt.Equal(update.StartsAt), "startsAt: expected %s, got %s", slot.StartsAt, update.StartsAt)
	assert.True(t, slot.EndsAt.Equal(update.EndsAt), "endsAt: expected %s, got %s", slot.EndsAt, update.EndsAt)
	assert.Equal(t, slot.IsValidated, update.IsValidated)
}

func TestNewSlotService_ReusesProvidedInstance(t *testing.T) {
	existing := &SlotService{}
	assert.Same(t, existing, NewSlotService(existing))
}

func TestSlotService_ConfirmSlot_EventCanceled(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	event.Status = constants.EVENT_STATUS_CANCELED
	require.NoError(t, s.eventRepository.Updates(&event))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
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
	assert.Equal(t, helpers.NewFieldError(constants.ERR_SLOT_INVALID_STARTS_AT, "startsAt"), err)
}

func TestSlotService_ConfirmSlot_InvalidEndsAt(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt.Add(time.Hour)}, slotEntity.Id, owner.Id)
	assert.Equal(t, helpers.NewFieldError(constants.ERR_SLOT_INVALID_ENDS_AT, "endsAt"), err)
}

func TestSlotService_ConfirmSlot_ReportsAllInvalidFields(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	_, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: event.StartsAt.Add(-time.Hour), EndsAt: slotEntity.EndsAt.Add(time.Hour)}, slotEntity.Id, owner.Id)
	assert.Equal(t, &helpers.ValidationError{Code: constants.ERR_INVALID_BODY, Fields: []helpers.FieldError{
		{Field: "startsAt", Code: constants.ERR_SLOT_INVALID_STARTS_AT.Err.Error()},
		{Field: "endsAt", Code: constants.ERR_SLOT_INVALID_ENDS_AT.Err.Error()},
	}}, err)
}

func TestSlotService_ConfirmSlot_Success(t *testing.T) {
	s := newTestSlotService(t)
	called := testutils.StubSMTPAwait(t, &s.mailService.SendMailFunc)

	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)

	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	resp, err := s.ConfirmSlot(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt}, slotEntity.Id, owner.Id)
	assert.NoError(t, err)
	assert.True(t, resp.IsValidated)

	var updatedEvent model.Event
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &updatedEvent))
	assert.Equal(t, constants.EVENT_STATUS_UPCOMING, updatedEvent.Status)

	// The validated slot is a new slot
	updates := readSlotUpdates(t, client)
	require.Len(t, updates, 1)
	assertSlotUpdate(t, constants.SLOT_ACTION_CREATED, model.Slot{Id: resp.Id, StartsAt: resp.StartsAt, EndsAt: resp.EndsAt, IsValidated: true}, updates[0])

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
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	err := s.RemoveValidatedSlot(slotEntity.Id, owner.Id)
	assert.NoError(t, err)
	// RemoveValidatedSlot also triggers LoadSlotsAsync, which reuses this
	// test's dedicated connection; wait for it before querying again.
	var updatedEvent model.Event
	s.WaitAsyncLoads()
	require.NoError(t, s.eventRepository.FindOneById(event.Id, &updatedEvent))
	assert.Equal(t, constants.EVENT_STATUS_IN_DECISION, updatedEvent.Status)

	updates := readSlotUpdates(t, client)
	require.Len(t, updates, 1)
	assertSlotUpdate(t, constants.SLOT_ACTION_DELETED, slotEntity, updates[0])

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
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	createTestAvailability(t, owner.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	createTestAvailability(t, participant.Id, event.Id, event.StartsAt.Add(15*time.Minute), event.StartsAt.Add(75*time.Minute))

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	require.Len(t, slots, 1)
	assert.True(t, slots[0].StartsAt.Equal(event.StartsAt.Add(15*time.Minute)))
	assert.True(t, slots[0].EndsAt.Equal(event.StartsAt.Add(time.Hour)))

	updates := readSlotUpdates(t, client)
	require.Len(t, updates, 1)
	assertSlotUpdate(t, constants.SLOT_ACTION_CREATED, slots[0], updates[0])
}

func TestSlotService_LoadSlots_UnchangedSlots_KeepIdsWithoutBroadcast(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)
	createTestAvailability(t, owner.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	createTestAvailability(t, participant.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	s.LoadSlots(event.Id)

	var slotsBefore []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsBefore))
	require.Len(t, slotsBefore, 1)
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	s.LoadSlots(event.Id)

	var slotsAfter []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsAfter))
	require.Len(t, slotsAfter, 1)
	assert.Equal(t, slotsBefore[0].Id, slotsAfter[0].Id)
	assert.Empty(t, client.Channel)
}

func TestSlotService_LoadSlots_AvailabilityUpdated_BroadcastsUpdatedWithSameId(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)
	createTestAvailability(t, owner.Id, event.Id, event.StartsAt, event.StartsAt.Add(2*time.Hour))
	participantAvailability := createTestAvailability(t, participant.Id, event.Id, event.StartsAt.Add(15*time.Minute), event.StartsAt.Add(75*time.Minute))
	s.LoadSlots(event.Id)

	var slotsBefore []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsBefore))
	require.Len(t, slotsBefore, 1)
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	participantAvailability.StartsAt = event.StartsAt.Add(30 * time.Minute)
	participantAvailability.EndsAt = event.StartsAt.Add(90 * time.Minute)
	require.NoError(t, repository.NewAvailabilityRepository(nil).Update(&participantAvailability))

	s.LoadSlots(event.Id)

	var slotsAfter []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsAfter))
	require.Len(t, slotsAfter, 1)
	assert.Equal(t, slotsBefore[0].Id, slotsAfter[0].Id)
	assert.True(t, slotsAfter[0].StartsAt.Equal(event.StartsAt.Add(30*time.Minute)))
	assert.True(t, slotsAfter[0].EndsAt.Equal(event.StartsAt.Add(90*time.Minute)))

	updates := readSlotUpdates(t, client)
	require.Len(t, updates, 1)
	assertSlotUpdate(t, constants.SLOT_ACTION_UPDATED, slotsAfter[0], updates[0])
}

func TestSlotService_LoadSlots_AvailabilityDeleted_DeletesSlotAndBroadcastsDeleted(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)
	createTestAvailability(t, owner.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	participantAvailability := createTestAvailability(t, participant.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	s.LoadSlots(event.Id)

	var slotsBefore []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsBefore))
	require.Len(t, slotsBefore, 1)
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	// Only one participant still has availabilities, so the slot must be removed
	require.NoError(t, repository.NewAvailabilityRepository(nil).DeleteById(&participantAvailability.Id))

	s.LoadSlots(event.Id)

	var slotsAfter []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slotsAfter))
	assert.Empty(t, slotsAfter)

	updates := readSlotUpdates(t, client)
	require.Len(t, updates, 1)
	assertSlotUpdate(t, constants.SLOT_ACTION_DELETED, slotsBefore[0], updates[0])
}

func TestSlotService_LoadSlots_ApplyChangesFails_NoBroadcast(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)
	createTestAvailability(t, owner.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	createTestAvailability(t, participant.Id, event.Id, event.StartsAt, event.StartsAt.Add(time.Hour))
	client := s.sseService.AddClient("client", owner.Id, event.Id, context.Background())

	testutils.MakeReadOnly(t)

	s.LoadSlots(event.Id)

	var slots []model.Slot
	require.NoError(t, s.slotRepository.FindByEventId(event.Id, &slots))
	assert.Empty(t, slots)
	assert.Empty(t, client.Channel)
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

func TestSlotService_LoadSlots_SlotsLookupFails_NoOp(t *testing.T) {
	s := newTestSlotService(t)
	owner := createTestAccount(t)
	participant := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id, participant.Id)

	s.slotRepository = repository.NewSlotRepository(testutils.ClosedDB(t))

	// Should not panic; the lookup error is logged and swallowed.
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
