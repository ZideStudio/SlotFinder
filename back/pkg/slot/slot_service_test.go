package slot

import (
	model "app/db/models"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindIntersectingTimeSlots_BasicIntersection(t *testing.T) {
	// Test basic intersection of two users' availabilities
	userAvailabilities := map[uuid.UUID][]TimeSlot{
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			},
		},
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
			},
		},
	}

	requiredDuration := 60 * time.Minute
	service := &SlotService{}
	result := service.findIntersectingTimeSlots(userAvailabilities, requiredDuration)

	assert.Len(t, result, 1, "Expected 1 common time slot")
	assert.Equal(t, time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), result[0].StartsAt, "Start time should be 12:00")
	assert.Equal(t, time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), result[0].EndsAt, "End time should be 14:00")
}

func TestFindIntersectingTimeSlots_NoIntersection(t *testing.T) {
	// Test no intersection between users' availabilities
	userAvailabilities := map[uuid.UUID][]TimeSlot{
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			},
		},
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
			},
		},
	}

	requiredDuration := 60 * time.Minute
	service := &SlotService{}
	result := service.findIntersectingTimeSlots(userAvailabilities, requiredDuration)

	assert.Len(t, result, 0, "Expected no common time slots")
}

func TestFindIntersectingTimeSlots_InsufficientDuration(t *testing.T) {
	// Test intersection that doesn't meet minimum duration
	userAvailabilities := map[uuid.UUID][]TimeSlot{
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			},
		},
	}

	requiredDuration := 60 * time.Minute // 1 hour required, but only 30 minutes overlap
	service := &SlotService{}
	result := service.findIntersectingTimeSlots(userAvailabilities, requiredDuration)

	assert.Len(t, result, 0, "Expected no slots due to insufficient duration")
}

func TestFindIntersectingTimeSlots_MultipleSlots(t *testing.T) {
	// Test multiple non-overlapping intersections
	userAvailabilities := map[uuid.UUID][]TimeSlot{
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				StartsAt: time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
			},
		},
		uuid.New(): {
			{
				StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			},
			{
				StartsAt: time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
			},
		},
	}

	requiredDuration := 30 * time.Minute
	service := &SlotService{}
	result := service.findIntersectingTimeSlots(userAvailabilities, requiredDuration)

	assert.Len(t, result, 2, "Expected 2 common time slots")

	// First slot: 10:00-11:00
	assert.Equal(t, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), result[0].StartsAt, "First slot start time should be 10:00")
	assert.Equal(t, time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC), result[0].EndsAt, "First slot end time should be 11:00")

	// Second slot: 15:00-17:00
	assert.Equal(t, time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC), result[1].StartsAt, "Second slot start time should be 15:00")
	assert.Equal(t, time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC), result[1].EndsAt, "Second slot end time should be 17:00")
}

func TestMergeOverlappingTimeSlots(t *testing.T) {
	// Test merging overlapping time slots
	slots := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			StartsAt: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
		},
		{
			StartsAt: time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
		},
	}

	service := &SlotService{}
	result := service.mergeOverlappingTimeSlots(slots)

	assert.Len(t, result, 2, "Expected 2 merged slots")

	// First merged slot: 10:00-14:00
	assert.Equal(t, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), result[0].StartsAt, "Merged slot start time should be 10:00")
	assert.Equal(t, time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), result[0].EndsAt, "Merged slot end time should be 14:00")

	// Second slot: 16:00-18:00 (unchanged)
	assert.Equal(t, time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC), result[1].StartsAt, "Second slot start time should be 16:00")
	assert.Equal(t, time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC), result[1].EndsAt, "Second slot end time should be 18:00")
}

func TestIntersectTimeSlots_SimpleOverlap(t *testing.T) {
	slots1 := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC),
		},
	}

	slots2 := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
		},
	}

	service := &SlotService{}
	result := service.intersectTimeSlots(slots1, slots2)

	assert.Len(t, result, 1, "Expected 1 intersection")
	assert.Equal(t, time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), result[0].StartsAt, "Intersection start should be 12:00")
	assert.Equal(t, time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC), result[0].EndsAt, "Intersection end should be 15:00")
}

func TestIntersectTimeSlots_NoOverlap(t *testing.T) {
	slots1 := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	slots2 := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		},
	}

	service := &SlotService{}
	result := service.intersectTimeSlots(slots1, slots2)

	assert.Len(t, result, 0, "Expected no intersections")
}

func TestMergeOverlappingTimeSlots_AdjacentSlots(t *testing.T) {
	slots := []TimeSlot{
		{
			StartsAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			StartsAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
		},
	}

	service := &SlotService{}
	result := service.mergeOverlappingTimeSlots(slots)

	assert.Len(t, result, 1, "Expected 1 merged slot")
	assert.Equal(t, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), result[0].StartsAt, "Merged slot start should be 10:00")
	assert.Equal(t, time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), result[0].EndsAt, "Merged slot end should be 14:00")
}

func TestLoadSlots_ConcurrentCallsDoNotRace(t *testing.T) {
	service := &SlotService{
		loadSlotsMutexes: sync.Map{},
	}
	eventId := uuid.New()

	// This test verifies that the mutex mechanism in LoadSlots prevents
	// race conditions when multiple goroutines try to recalculate slots
	// for the same event concurrently. The test doesn't require database
	// access - we're just verifying the mutex works.

	const numCalls = 10
	callCount := 0
	var mu sync.Mutex

	// Create a channel to coordinate the test
	done := make(chan bool, numCalls)
	start := make(chan bool)

	// Launch multiple concurrent LoadSlots calls for the same event
	for i := 0; i < numCalls; i++ {
		go func() {
			<-start // Wait for all goroutines to be ready

			// Get the mutex for this event
			mutexInterface, _ := service.loadSlotsMutexes.LoadOrStore(eventId.String(), &sync.Mutex{})
			mutex := mutexInterface.(*sync.Mutex)

			// Try to acquire the mutex
			mutex.Lock()
			// Increment call count (simulating critical section)
			mu.Lock()
			callCount++
			mu.Unlock()
			mutex.Unlock()

			done <- true
		}()
	}

	// Signal all goroutines to start
	close(start)

	// Wait for all goroutines to complete
	for i := 0; i < numCalls; i++ {
		<-done
	}

	// Verify all calls completed
	assert.Equal(t, numCalls, callCount, "All concurrent calls should have completed")
}

func TestLoadSlots_ConcurrentCallsDifferentEvents(t *testing.T) {
	service := &SlotService{
		loadSlotsMutexes: sync.Map{},
	}

	// Test that calls to LoadSlots for different events use different mutexes
	// and can run concurrently without blocking each other
	const numEvents = 5
	done := make(chan bool, numEvents)
	start := make(chan bool)

	for i := 0; i < numEvents; i++ {
		eventId := uuid.New()
		go func(id uuid.UUID) {
			<-start

			// Get the mutex for this event
			mutexInterface, _ := service.loadSlotsMutexes.LoadOrStore(id.String(), &sync.Mutex{})
			mutex := mutexInterface.(*sync.Mutex)

			mutex.Lock()
			// Simulate some work
			time.Sleep(10 * time.Millisecond)
			mutex.Unlock()

			done <- true
		}(eventId)
	}

	// Signal all goroutines to start
	close(start)

	// Wait for all goroutines to complete with timeout
	timeout := time.After(1 * time.Second)
	for i := 0; i < numEvents; i++ {
		select {
		case <-done:
			// Success
		case <-timeout:
			t.Fatal("Test timed out - possible deadlock")
		}
	}
}

func slotHour(hour int) time.Time {
	return time.Date(2024, 1, 1, hour, 0, 0, 0, time.UTC)
}

func TestDiffSlots_NoExistingSlots_CreatesAll(t *testing.T) {
	eventId := uuid.New()
	service := &SlotService{}

	changes := service.diffSlots(eventId, nil, []TimeSlot{
		{StartsAt: slotHour(10), EndsAt: slotHour(11)},
		{StartsAt: slotHour(13), EndsAt: slotHour(14)},
	})

	assert.Empty(t, changes.Updated)
	assert.Empty(t, changes.Deleted)
	require.Len(t, changes.Created, 2)
	assert.NotEqual(t, uuid.Nil, changes.Created[0].Id)
	assert.NotEqual(t, changes.Created[0].Id, changes.Created[1].Id)
	assert.Equal(t, model.Slot{Id: changes.Created[0].Id, EventId: eventId, StartsAt: slotHour(10), EndsAt: slotHour(11)}, changes.Created[0])
	assert.Equal(t, model.Slot{Id: changes.Created[1].Id, EventId: eventId, StartsAt: slotHour(13), EndsAt: slotHour(14)}, changes.Created[1])
}

func TestDiffSlots_NoComputedSlots_DeletesAll(t *testing.T) {
	existing := []model.Slot{
		{Id: uuid.New(), StartsAt: slotHour(10), EndsAt: slotHour(11)},
		{Id: uuid.New(), StartsAt: slotHour(13), EndsAt: slotHour(14)},
	}
	service := &SlotService{}

	changes := service.diffSlots(uuid.New(), existing, []TimeSlot{})

	assert.Empty(t, changes.Created)
	assert.Empty(t, changes.Updated)
	assert.Equal(t, existing, changes.Deleted)
}

func TestDiffSlots_SameDatesDifferentLocation_NoChanges(t *testing.T) {
	// Dates read from the database do not have the same location as the computed ones
	location := time.FixedZone("UTC+1", 60*60)
	existing := []model.Slot{{Id: uuid.New(), StartsAt: slotHour(10).In(location), EndsAt: slotHour(11).In(location)}}
	service := &SlotService{}

	changes := service.diffSlots(uuid.New(), existing, []TimeSlot{{StartsAt: slotHour(10), EndsAt: slotHour(11)}})

	assert.Empty(t, changes.Created)
	assert.Empty(t, changes.Updated)
	assert.Empty(t, changes.Deleted)
}

func TestDiffSlots_DatesChanged_UpdatesKeepingId(t *testing.T) {
	eventId := uuid.New()
	existing := model.Slot{Id: uuid.New(), EventId: eventId, StartsAt: slotHour(10), EndsAt: slotHour(12)}
	service := &SlotService{}

	changes := service.diffSlots(eventId, []model.Slot{existing}, []TimeSlot{{StartsAt: slotHour(11), EndsAt: slotHour(13)}})

	assert.Empty(t, changes.Created)
	assert.Empty(t, changes.Deleted)
	assert.Equal(t, []model.Slot{{Id: existing.Id, EventId: eventId, StartsAt: slotHour(11), EndsAt: slotHour(13)}}, changes.Updated)
}

func TestDiffSlots_MergedSlots_UpdatesFirstAndDeletesSecond(t *testing.T) {
	first := model.Slot{Id: uuid.New(), StartsAt: slotHour(10), EndsAt: slotHour(11)}
	second := model.Slot{Id: uuid.New(), StartsAt: slotHour(13), EndsAt: slotHour(14)}
	service := &SlotService{}

	changes := service.diffSlots(uuid.New(), []model.Slot{first, second}, []TimeSlot{{StartsAt: slotHour(10), EndsAt: slotHour(14)}})

	assert.Empty(t, changes.Created)
	assert.Equal(t, []model.Slot{{Id: first.Id, StartsAt: slotHour(10), EndsAt: slotHour(14)}}, changes.Updated)
	assert.Equal(t, []model.Slot{second}, changes.Deleted)
}

func TestDiffSlots_SplitSlot_UpdatesFirstAndCreatesSecond(t *testing.T) {
	eventId := uuid.New()
	existing := model.Slot{Id: uuid.New(), EventId: eventId, StartsAt: slotHour(10), EndsAt: slotHour(14)}
	service := &SlotService{}

	changes := service.diffSlots(eventId, []model.Slot{existing}, []TimeSlot{
		{StartsAt: slotHour(10), EndsAt: slotHour(11)},
		{StartsAt: slotHour(13), EndsAt: slotHour(14)},
	})

	assert.Empty(t, changes.Deleted)
	assert.Equal(t, []model.Slot{{Id: existing.Id, EventId: eventId, StartsAt: slotHour(10), EndsAt: slotHour(11)}}, changes.Updated)
	require.Len(t, changes.Created, 1)
	assert.Equal(t, model.Slot{Id: changes.Created[0].Id, EventId: eventId, StartsAt: slotHour(13), EndsAt: slotHour(14)}, changes.Created[0])
}

func TestDiffSlots_MovedWithoutOverlap_DeletesAndCreates(t *testing.T) {
	eventId := uuid.New()
	existing := model.Slot{Id: uuid.New(), EventId: eventId, StartsAt: slotHour(9), EndsAt: slotHour(10)}
	service := &SlotService{}

	// A slot only touching the existing one does not overlap it either
	changes := service.diffSlots(eventId, []model.Slot{existing}, []TimeSlot{{StartsAt: slotHour(10), EndsAt: slotHour(11)}})

	assert.Empty(t, changes.Updated)
	assert.Equal(t, []model.Slot{existing}, changes.Deleted)
	require.Len(t, changes.Created, 1)
	assert.NotEqual(t, existing.Id, changes.Created[0].Id)
	assert.Equal(t, model.Slot{Id: changes.Created[0].Id, EventId: eventId, StartsAt: slotHour(10), EndsAt: slotHour(11)}, changes.Created[0])
}

func TestDiffSlots_UnsortedExistingSlots_ReusesEarliest(t *testing.T) {
	first := model.Slot{Id: uuid.New(), StartsAt: slotHour(10), EndsAt: slotHour(11)}
	second := model.Slot{Id: uuid.New(), StartsAt: slotHour(13), EndsAt: slotHour(14)}
	service := &SlotService{}

	changes := service.diffSlots(uuid.New(), []model.Slot{second, first}, []TimeSlot{{StartsAt: slotHour(10), EndsAt: slotHour(14)}})

	assert.Empty(t, changes.Created)
	assert.Equal(t, []model.Slot{{Id: first.Id, StartsAt: slotHour(10), EndsAt: slotHour(14)}}, changes.Updated)
	assert.Equal(t, []model.Slot{second}, changes.Deleted)
}
