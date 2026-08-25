package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	appdb "app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/slot"
	"app/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestAvailabilityService(t *testing.T) *AvailabilityService {
	t.Helper()
	testutils.TestDB(t)
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
	require.NoError(t, appdb.GetDB().Create(&model.Account{Id: id, Username: &username}).Error)
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

// createEndedEventWithAccess creates an event whose EndsAt is already in the
// past (and status not yet FINISHED), so CheckAndAutoUpdateStatus attempts
// to auto-transition it via eventRepository.Updates.
func createEndedEventWithAccess(t *testing.T, db *repository.EventRepository, ownerId uuid.UUID) model.Event {
	t.Helper()
	start := time.Now().UTC().Add(-4 * time.Hour).Truncate(5 * time.Minute)
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Ended Event",
		Duration: 30,
		StartsAt: start,
		EndsAt:   start.Add(time.Hour),
		OwnerId:  ownerId,
		Status:   constants.EVENT_STATUS_IN_DECISION,
	}
	require.NoError(t, db.Create(&event))

	accountEventRepo := repository.NewAccountEventRepository(nil)
	require.NoError(t, accountEventRepo.Create(&model.AccountEvent{AccountId: ownerId, EventId: event.Id}))

	var found model.Event
	require.NoError(t, db.FindOneById(event.Id, &found))
	return found
}

func TestAvailabilityService_Create_AutoFinishUpdateError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEndedEventWithAccess(t, s.eventRepository, owner)
	testutils.MakeReadOnly(t)

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
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
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.availabilityRepository.FindByEventId(event.Id, &availabilities); err != nil {
			return false
		}
		return len(availabilities) == 1
	})
	assert.Len(t, availabilities, 1, "overlapping availabilities should have been merged")
}

// TestAvailabilityService_Create_MergesOverlapping_ExtendsEnd covers the
// merge branch where the *existing* availability's EndsAt is later than the
// one being created, so it's the existing end that wins.
func TestAvailabilityService_Create_MergesOverlapping_ExtendsEnd(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(2 * time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	// Contained within the first (0-2h): overlaps, but its own EndsAt (1h)
	// is earlier than the existing one's (2h) -> existing end should win.
	dto, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)
	assert.True(t, dto.EndsAt.Equal(event.StartsAt.Add(2*time.Hour)))
}

func TestAvailabilityService_Create_FindOverlappingError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	s.availabilityRepository = repository.NewAvailabilityRepository(testutils.ClosedDB(t))

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestAvailabilityService_Create_RepositoryCreateError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	testutils.MakeReadOnly(t)

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestAvailabilityService_Create_MergePath_DeleteByIdsError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	// Overlaps with the first -> enters the merge path, whose DeleteByIds call fails.
	_, err = s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(90 * time.Minute)}, event.Id, claims)
	assert.Error(t, err)
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

func TestAvailabilityService_Update_FindOneByIdError(t *testing.T) {
	s := newTestAvailabilityService(t)
	s.availabilityRepository = repository.NewAvailabilityRepository(testutils.ClosedDB(t))

	_, err := s.Update(&AvailabilityUpdateDto{}, uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
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

func TestAvailabilityService_Update_NoMergeCase_RepositoryUpdateError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	newEnd := event.StartsAt.Add(90 * time.Minute)
	_, err = s.Update(&AvailabilityUpdateDto{EndsAt: &newEnd}, created.Id, claims)
	assert.Error(t, err)
}

// Covers the merge branch where the existing overlapping availability's
// EndsAt is later than the one being updated, so the existing end wins.
func TestAvailabilityService_Update_MergesOverlapping_ExtendsEnd(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	// Wide, will be "existing" once the second one is updated to overlap it.
	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(2 * time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	// Initially non-overlapping.
	second, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(3 * time.Hour), EndsAt: event.StartsAt.Add(3*time.Hour + 30*time.Minute)}, event.Id, claims)
	require.NoError(t, err)

	// Move it to be fully contained within the wide one's range -> overlaps,
	// and its own EndsAt is earlier than the existing (wide) one's -> existing wins.
	newStart := event.StartsAt.Add(30 * time.Minute)
	newEnd := event.StartsAt.Add(time.Hour)
	updated, err := s.Update(&AvailabilityUpdateDto{StartsAt: &newStart, EndsAt: &newEnd}, second.Id, claims)
	require.NoError(t, err)
	assert.True(t, updated.EndsAt.Equal(event.StartsAt.Add(2*time.Hour)))
}

// The merged range must still fit the event's window. Here the event is
// shrunk after both availabilities were created, so the merged EndsAt
// (from the untouched "wide" one) now falls outside it.
func TestAvailabilityService_Update_MergePath_RevalidationFails(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	// Wide, stays untouched and will lend its EndsAt to the merged result.
	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(3 * time.Hour), EndsAt: event.StartsAt.Add(3*time.Hour + 50*time.Minute)}, event.Id, claims)
	require.NoError(t, err)

	// Initially non-overlapping.
	second, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}, event.Id, claims)
	require.NoError(t, err)

	// Shrink the event so the wide availability's EndsAt (3h50) no longer fits,
	// while still leaving room for the moved availability's own range to validate.
	event.EndsAt = event.StartsAt.Add(3*time.Hour + 30*time.Minute)
	require.NoError(t, s.eventRepository.Updates(&event))

	// Move it to overlap the wide one; its own new range (2h55-3h05) still fits
	// the shrunk event, so only the post-merge revalidation can catch it.
	newStart := event.StartsAt.Add(2*time.Hour + 55*time.Minute)
	newEnd := event.StartsAt.Add(3*time.Hour + 5*time.Minute)
	_, err = s.Update(&AvailabilityUpdateDto{StartsAt: &newStart, EndsAt: &newEnd}, second.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_END_AFTER_EVENT.Err)
}

func TestAvailabilityService_Update_MergePath_DeleteByIdsError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	_, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)
	second, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt.Add(2 * time.Hour), EndsAt: event.StartsAt.Add(3 * time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	// Extend the second so it now overlaps the first -> merge path, whose DeleteByIds call fails.
	newStart := event.StartsAt.Add(30 * time.Minute)
	_, err = s.Update(&AvailabilityUpdateDto{StartsAt: &newStart}, second.Id, claims)
	assert.Error(t, err)
}

func TestAvailabilityService_Delete_NotFound(t *testing.T) {
	s := newTestAvailabilityService(t)
	err := s.Delete(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.ErrorIs(t, err, constants.ERR_AVAILABILITY_NOT_FOUND.Err)
}

func TestAvailabilityService_Delete_FindOneByIdError(t *testing.T) {
	s := newTestAvailabilityService(t)
	s.availabilityRepository = repository.NewAvailabilityRepository(testutils.ClosedDB(t))

	err := s.Delete(uuid.New(), &guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
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
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.availabilityRepository.FindByEventId(event.Id, &availabilities); err != nil {
			return false
		}
		return len(availabilities) == 0
	})
	assert.Len(t, availabilities, 0)
}

// TestAvailabilityService_Delete_EventAccessDenied covers the case where the
// caller owns the availability but was later removed from the event.
func TestAvailabilityService_Delete_EventAccessDenied(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	require.NoError(t, appdb.GetDB().Where("account_id = ? AND event_id = ?", owner, event.Id).Delete(&model.AccountEvent{}).Error)

	err = s.Delete(created.Id, claims)
	assert.ErrorIs(t, err, constants.ERR_EVENT_ACCESS_DENIED.Err)
}

func TestAvailabilityService_Delete_AutoFinishUpdateError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEndedEventWithAccess(t, s.eventRepository, owner)

	// Insert directly (bypassing Create, which would itself reject an ended event).
	availability := model.Availability{Id: uuid.New(), StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour), AccountId: owner, EventId: event.Id}
	require.NoError(t, s.availabilityRepository.Create(&availability))

	testutils.MakeReadOnly(t)

	err := s.Delete(availability.Id, &guard.Claims{Id: owner})
	assert.Error(t, err)
}

func TestAvailabilityService_Delete_RepositoryDeleteError(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	err = s.Delete(created.Id, claims)
	assert.Error(t, err)
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
	testutils.AwaitAsyncDBWorkUntil(t, 2*time.Second, func() bool {
		if err := s.availabilityRepository.FindByEventId(event.Id, &availabilities); err != nil {
			return false
		}
		return len(availabilities) == 1
	})
	assert.Len(t, availabilities, 1, "overlapping availabilities should have been merged")
}
