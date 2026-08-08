package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	appdb "app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/slot"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestAvailabilityService(t *testing.T) *AvailabilityService {
	t.Helper()
	return &AvailabilityService{
		slotService:            slot.NewSlotService(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		eventRepository:        repository.NewEventRepository(nil),
	}
}

// alignedNow returns a time truncated to a 5-minute boundary in UTC, a
// minute in the future so it's safely "not started yet".
func alignedNow(t *testing.T) time.Time {
	t.Helper()
	return time.Now().UTC().Add(time.Hour).Truncate(5 * time.Minute)
}

func createEventWithAccess(t *testing.T, db *repository.EventRepository, ownerId uuid.UUID) model.Event {
	t.Helper()
	start := alignedNow(t)
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Test Event",
		Duration: 30,
		StartsAt: start,
		EndsAt:   start.Add(4 * time.Hour),
		OwnerId:  ownerId,
		Status:   constants.EVENT_STATUS_IN_DECISION,
	}
	require.NoError(t, db.Create(&event))

	accountEventRepo := repository.NewAccountEventRepository(nil)
	require.NoError(t, accountEventRepo.Create(&model.AccountEvent{AccountId: ownerId, EventId: event.Id}))

	// Re-fetch so AccountEvents/Owner are populated (HasUserAccess needs AccountEvents).
	var found model.Event
	require.NoError(t, db.FindOneById(event.Id, &found))
	return found
}

func createTestUser(t *testing.T) uuid.UUID {
	t.Helper()
	id := uuid.New()
	username := "avail-" + uuid.NewString()
	require.NoError(t, appdb.GetDB().Create(&model.Account{Id: id, UserName: &username}).Error)
	return id
}

func TestNewAvailabilityService_Nil_BuildsRealDependencies(t *testing.T) {
	s := NewAvailabilityService(nil)
	assert.NotNil(t, s.slotService)
	assert.NotNil(t, s.availabilityRepository)
	assert.NotNil(t, s.eventRepository)
}

func TestNewAvailabilityService_NonNil_ReturnsSameInstance(t *testing.T) {
	existing := &AvailabilityService{}
	s := NewAvailabilityService(existing)
	assert.Same(t, existing, s)
}

// createUpcomingEvent creates an event with status UPCOMING (not IN_DECISION) and a future EndsAt.
func createUpcomingEvent(t *testing.T, db *repository.EventRepository, ownerId uuid.UUID) model.Event {
	t.Helper()
	start := alignedNow(t)
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Upcoming Event",
		Duration: 30,
		StartsAt: start,
		EndsAt:   start.Add(4 * time.Hour),
		OwnerId:  ownerId,
		Status:   constants.EVENT_STATUS_UPCOMING,
	}
	require.NoError(t, db.Create(&event))

	accountEventRepo := repository.NewAccountEventRepository(nil)
	require.NoError(t, accountEventRepo.Create(&model.AccountEvent{AccountId: ownerId, EventId: event.Id}))

	var found model.Event
	require.NoError(t, db.FindOneById(event.Id, &found))
	return found
}

func TestAvailabilityService_Create_EventNotFound(t *testing.T) {
	s := newTestAvailabilityService(t)
	user := &guard.Claims{Id: uuid.New()}

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: alignedNow(t), EndsAt: alignedNow(t).Add(time.Hour)}, uuid.New(), user)
	assert.ErrorIs(t, err, constants.ERR_EVENT_NOT_FOUND.Err)
}

func TestAvailabilityService_Create_AccessDenied(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)

	stranger := &guard.Claims{Id: uuid.New()}
	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, stranger)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ACCESS_DENIED.Err)
}

func TestAvailabilityService_Create_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)

	dto, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}, event.Id, &guard.Claims{Id: owner})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, dto.Id)
}

func TestAvailabilityService_Create_MergesOverlapping(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	// Overlaps with the first availability -> should merge into one.
	dto, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(90 * time.Minute)}, event.Id, claims)
	require.NoError(t, err)
	assert.True(t, dto.EndsAt.Equal(event.StartsAt.Add(90*time.Minute)))

	var availabilities []model.Availability
	require.NoError(t, s.availabilityRepository.FindByEventId(event.Id, &availabilities))
	assert.Len(t, availabilities, 1, "overlapping availabilities should have been merged")
}

func TestAvailabilityService_Create_InvalidTimes(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	// End before start -> validateAvailabilityTimes rejects it.
	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(time.Hour), EndsAt: event.StartsAt}, event.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_EVENT_START_AFTER_END.Err)
}

func TestAvailabilityService_Update_InvalidTimes(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	// New end before existing start -> validateAvailabilityTimes rejects it.
	newEnd := event.StartsAt.Add(-time.Hour)
	_, err = s.Update(&AvailabilityUpdateDto{EndsAt: &newEnd}, created.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_EVENT_START_AFTER_END.Err)
}

func TestAvailabilityService_Update_NotFound(t *testing.T) {
	s := newTestAvailabilityService(t)
	_, err := s.Update(&AvailabilityUpdateDto{}, uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_NOT_FOUND.Err)
}

func TestAvailabilityService_Update_AccessDenied(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	stranger := &guard.Claims{Id: uuid.New()}
	_, err = s.Update(&AvailabilityUpdateDto{}, created.Id, stranger)
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_ACCESS_DENIED.Err)
}

func TestAvailabilityService_Update_NoFieldsProvided_NoOp(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	updated, err := s.Update(&AvailabilityUpdateDto{}, created.Id, claims)
	assert.NoError(t, err)
	assert.True(t, updated.EndsAt.Equal(created.EndsAt))
}

func TestAvailabilityService_Update_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	newEnd := event.StartsAt.Add(90 * time.Minute)
	updated, err := s.Update(&AvailabilityUpdateDto{EndsAt: &newEnd}, created.Id, claims)
	assert.NoError(t, err)
	assert.True(t, updated.EndsAt.Equal(newEnd))
}

func TestAvailabilityService_Delete_NotFound(t *testing.T) {
	s := newTestAvailabilityService(t)
	err := s.Delete(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_NOT_FOUND.Err)
}

func TestAvailabilityService_Delete_AccessDenied(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	stranger := &guard.Claims{Id: uuid.New()}
	err = s.Delete(created.Id, stranger)
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_ACCESS_DENIED.Err)
}

func TestAvailabilityService_Delete_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	err = s.Delete(created.Id, claims)
	assert.NoError(t, err)

	var availabilities []model.Availability
	require.NoError(t, s.availabilityRepository.FindByEventId(event.Id, &availabilities))
	assert.Len(t, availabilities, 0)
}

func TestAvailabilityService_Update_EventEnded(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createUpcomingEvent(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	// Bypass Create (which would itself reject an UPCOMING event) to isolate Update's own check.
	availability := model.Availability{
		Id:        uuid.New(),
		StartsAt:  event.StartsAt,
		EndsAt:    event.StartsAt.Add(time.Hour),
		AccountId: owner,
		EventId:   event.Id,
	}
	require.NoError(t, s.availabilityRepository.Create(&availability))

	newEnd := event.StartsAt.Add(90 * time.Minute)
	_, err := s.Update(&AvailabilityUpdateDto{EndsAt: &newEnd}, availability.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
}

func TestAvailabilityService_Delete_EventEnded(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createUpcomingEvent(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	// Bypass Create (which would itself reject an UPCOMING event) to isolate Delete's own check.
	availability := model.Availability{
		Id:        uuid.New(),
		StartsAt:  event.StartsAt,
		EndsAt:    event.StartsAt.Add(time.Hour),
		AccountId: owner,
		EventId:   event.Id,
	}
	require.NoError(t, s.availabilityRepository.Create(&availability))

	err := s.Delete(availability.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ENDED.Err)
}

func TestAvailabilityService_Update_MergesOverlapping(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	first, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	second, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(2 * time.Hour), EndsAt: event.StartsAt.Add(3 * time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	// Extend the second availability so it now overlaps the first -> merge.
	newStart := event.StartsAt.Add(30 * time.Minute)
	updated, err := s.Update(&AvailabilityUpdateDto{StartsAt: &newStart}, second.Id, claims)
	require.NoError(t, err)
	assert.True(t, updated.StartsAt.Equal(first.StartsAt))

	var availabilities []model.Availability
	require.NoError(t, s.availabilityRepository.FindByEventId(event.Id, &availabilities))
	assert.Len(t, availabilities, 1, "overlapping availabilities should have been merged")
}
