package slot

import (
	"app/commons/constants"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/sse"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SlotService struct {
	slotRepository         *repository.SlotRepository
	eventRepository        *repository.EventRepository
	availabilityRepository *repository.AvailabilityRepository
	accountEventRepository *repository.AccountEventRepository
	sseService             *sse.SSEService
	loadSlotsMutexes       sync.Map // Map of eventId to *sync.Mutex for preventing concurrent LoadSlots
}

func NewSlotService(service *SlotService) *SlotService {
	if service != nil {
		return service
	}

	return &SlotService{
		slotRepository:         &repository.SlotRepository{},
		eventRepository:        &repository.EventRepository{},
		availabilityRepository: &repository.AvailabilityRepository{},
		accountEventRepository: &repository.AccountEventRepository{},
		sseService:             sse.GetSSEService(),
		loadSlotsMutexes:       sync.Map{},
	}
}

// Time interval
type TimeSlot struct {
	StartsAt time.Time
	EndsAt   time.Time
}

func (s *SlotService) ConfirmSlot(dto ConfirmSlotDto, slotId uuid.UUID, userId uuid.UUID) (model.Slot, error) {
	var selectedSlot model.Slot
	if err := s.slotRepository.FindOneById(slotId, &selectedSlot); err != nil {
		return model.Slot{}, constants.ERR_SLOT_NOT_FOUND.Err
	}

	// Check if user is admin of the event
	if !selectedSlot.Event.IsOwner(&userId) {
		return model.Slot{}, constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	// Check if event is locked
	if selectedSlot.Event.IsLocked() {
		return model.Slot{}, constants.ERR_EVENT_ENDED.Err
	}

	// Check if dto StartsAt is equals or after selectedSlot.StartsAt and before selectedSlot.EndsAt
	if dto.StartsAt.Before(selectedSlot.StartsAt) || !dto.StartsAt.Before(selectedSlot.EndsAt) {
		return model.Slot{}, constants.ERR_SLOT_INVALID_STARTS_AT.Err
	}
	// Check if dto EndsAt is after dto.StartsAt and before or equals selectedSlot.EndsAt
	if !dto.EndsAt.After(dto.StartsAt) || dto.EndsAt.After(selectedSlot.EndsAt) {
		return model.Slot{}, constants.ERR_SLOT_INVALID_ENDS_AT.Err
	}

	// Create a new validated slot from the selected slot
	slot := model.Slot{
		Id:          uuid.New(),
		EventId:     selectedSlot.EventId,
		StartsAt:    dto.StartsAt,
		EndsAt:      dto.EndsAt,
		IsValidated: true,
	}
	if err := s.slotRepository.Create(&slot); err != nil {
		return model.Slot{}, err
	}

	// Update event status
	event := model.Event{
		Id:     selectedSlot.EventId,
		Status: constants.EVENT_STATUS_UPCOMING,
	}
	if err := s.eventRepository.Updates(&event); err != nil {
		return model.Slot{}, err
	}

	return slot, nil
}

// Recalculates and recreates all slots for an event
func (s *SlotService) LoadSlots(eventId uuid.UUID) {
	// Acquire per-event mutex to prevent concurrent slot recalculations for the same event
	mutexInterface, _ := s.loadSlotsMutexes.LoadOrStore(eventId.String(), &sync.Mutex{})
	mutex := mutexInterface.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()

	log.Debug().Str("eventId", eventId.String()).Msg("Starting slot recalculation")

	// Get event
	var event model.Event
	if err := s.eventRepository.FindOneById(eventId, &event); err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to get event for slot calculation")
		return
	}

	// If event is finished, do not recalculate slots
	if event.IsLocked() {
		log.Debug().Str("eventId", eventId.String()).Msg("Event is locked, skipping slot recalculation")
		return
	}

	// Delete existing slots for this event
	if err := s.slotRepository.DeleteByEventId(eventId); err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to delete existing slots")
		return
	}

	accountEvents := event.AccountEvents
	if len(accountEvents) < 2 {
		log.Debug().Str("eventId", eventId.String()).Msg("Not enough participants to calculate slots")
		return
	}

	// Get all availabilities for this event
	var availabilities []model.Availability
	if err := s.availabilityRepository.FindByEventId(eventId, &availabilities); err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to get availabilities")
		return
	}

	// Get all active user IDs and their availabilities
	userAvailabilities := make(map[uuid.UUID][]TimeSlot)
	for _, availability := range availabilities {
		userAvailabilities[availability.AccountId] = append(
			userAvailabilities[availability.AccountId],
			TimeSlot{
				StartsAt: availability.StartsAt,
				EndsAt:   availability.EndsAt,
			},
		)
	}

	// Find common available time slots
	commonSlots := s.findIntersectingTimeSlots(userAvailabilities, time.Duration(event.Duration)*time.Minute)
	if len(commonSlots) == 0 {
		log.Debug().Str("eventId", eventId.String()).Msg("No common available slots found")
		return
	}

	// Create new slots in database
	slots := make([]model.Slot, 0, len(commonSlots))
	for _, slot := range commonSlots {
		newSlot := model.Slot{
			Id:          uuid.New(),
			EventId:     eventId,
			StartsAt:    slot.StartsAt,
			EndsAt:      slot.EndsAt,
			IsValidated: false,
		}

		if err := s.slotRepository.Create(&newSlot); err != nil {
			log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to create slot")
		}

		slots = append(slots, newSlot)
	}

	log.Debug().Str("eventId", eventId.String()).Int("slotsCreated", len(commonSlots)).Msg("Slot recalculation completed")

	// Send new slots to all participants via SSE
	s.sseService.BroadcastSlotsUpdate(eventId, slots)
}

// Finds time slots where all users are available
func (s *SlotService) findIntersectingTimeSlots(userAvailabilities map[uuid.UUID][]TimeSlot, requiredDuration time.Duration) []TimeSlot {
	// Get all active user IDs
	allUserIds := make([]uuid.UUID, 0, len(userAvailabilities))
	for userID := range userAvailabilities {
		allUserIds = append(allUserIds, userID)
	}
	if len(allUserIds) < 2 {
		return []TimeSlot{}
	}

	// Intersect with each subsequent user's availabilities
	commonSlots := userAvailabilities[allUserIds[0]]
	for i := 1; i < len(allUserIds); i++ {
		userSlots := userAvailabilities[allUserIds[i]]
		commonSlots = s.intersectTimeSlots(commonSlots, userSlots)
		if len(commonSlots) == 0 {
			break
		}
	}

	// Filter slots that meet the minimum duration requirement
	validSlots := []TimeSlot{}
	for _, slot := range commonSlots {
		if slot.EndsAt.Sub(slot.StartsAt) >= requiredDuration {
			validSlots = append(validSlots, slot)
		}
	}

	return validSlots
}

// Finds the intersection of two sets of time slots
func (s *SlotService) intersectTimeSlots(slots1, slots2 []TimeSlot) []TimeSlot {
	var intersections []TimeSlot

	for _, slot1 := range slots1 {
		for _, slot2 := range slots2 {
			// Find overlap
			start := slot1.StartsAt
			if slot2.StartsAt.After(start) {
				start = slot2.StartsAt
			}

			end := slot1.EndsAt
			if slot2.EndsAt.Before(end) {
				end = slot2.EndsAt
			}

			// Overlap
			if start.Before(end) {
				intersections = append(intersections, TimeSlot{
					StartsAt: start,
					EndsAt:   end,
				})
			}
		}
	}

	return s.mergeOverlappingTimeSlots(intersections)
}

// Merges overlapping time slots
func (s *SlotService) mergeOverlappingTimeSlots(slots []TimeSlot) []TimeSlot {
	if len(slots) <= 1 {
		return slots
	}

	// Sort slots by start time
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartsAt.Before(slots[j].StartsAt)
	})

	merged := []TimeSlot{slots[0]}

	for i := 1; i < len(slots); i++ {
		lastMerged := &merged[len(merged)-1]
		current := slots[i]

		// If current slot starts before or at the end of the last merged slot, merge them
		if (current.StartsAt.Before(lastMerged.EndsAt) || current.StartsAt.Equal(lastMerged.EndsAt)) && current.EndsAt.After(lastMerged.EndsAt) {
			lastMerged.EndsAt = current.EndsAt
		} else if !current.StartsAt.Before(lastMerged.EndsAt) && !current.StartsAt.Equal(lastMerged.EndsAt) {
			// No overlap, add a new slot
			merged = append(merged, current)
		}
	}

	return merged
}
