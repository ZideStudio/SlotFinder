package slot

import (
	"app/commons/constants"
	"app/commons/helpers"
	"app/config"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
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
	mailService            *mail.MailService
	config                 *config.Config
	loadSlotsMutexes       sync.Map       // Map of eventId to *sync.Mutex for preventing concurrent LoadSlots
	asyncLoads             sync.WaitGroup // Tracks in-flight LoadSlotsAsync calls
}

func NewSlotService(service *SlotService) *SlotService {
	if service != nil {
		return service
	}

	return &SlotService{
		slotRepository:         repository.NewSlotRepository(nil),
		eventRepository:        repository.NewEventRepository(nil),
		availabilityRepository: repository.NewAvailabilityRepository(nil),
		accountEventRepository: repository.NewAccountEventRepository(nil),
		sseService:             sse.GetSSEService(),
		mailService:            mail.NewMailService(nil),
		loadSlotsMutexes:       sync.Map{},
		config:                 config.GetConfig(),
	}
}

// Time interval
type TimeSlot struct {
	StartsAt time.Time
	EndsAt   time.Time
}

// Slots touched by a recalculation
type SlotChanges struct {
	Created []model.Slot
	Updated []model.Slot
	Deleted []model.Slot
}

func (s *SlotService) ConfirmSlot(dto ConfirmSlotDto, slotId uuid.UUID, userId uuid.UUID) (SlotResponseDto, error) {
	var selectedSlot model.Slot
	if err := s.slotRepository.FindOneById(slotId, &selectedSlot); err != nil {
		return SlotResponseDto{}, constants.ERR_SLOT_NOT_FOUND.Err
	}

	// Check if user is admin of the event
	if !selectedSlot.Event.IsOwner(&userId) {
		return SlotResponseDto{}, constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	// Check if event is locked
	if hasStatus, err := selectedSlot.Event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, &[]constants.EventStatus{constants.EVENT_STATUS_IN_DECISION}); !hasStatus || err != nil {
		if err != nil {
			return SlotResponseDto{}, err
		}
		return SlotResponseDto{}, constants.ERR_EVENT_ENDED.Err
	}

	var fieldErrors helpers.FieldErrors
	// Check if dto StartsAt is equals or after selectedSlot.StartsAt and before selectedSlot.EndsAt
	if dto.StartsAt.Before(selectedSlot.StartsAt) || !dto.StartsAt.Before(selectedSlot.EndsAt) {
		fieldErrors.Add(constants.ERR_SLOT_INVALID_STARTS_AT, "startsAt")
	}
	// Check if dto EndsAt is after dto.StartsAt and before or equals selectedSlot.EndsAt
	if !dto.EndsAt.After(dto.StartsAt) || dto.EndsAt.After(selectedSlot.EndsAt) {
		fieldErrors.Add(constants.ERR_SLOT_INVALID_ENDS_AT, "endsAt")
	}
	if err := fieldErrors.Err(); err != nil {
		return SlotResponseDto{}, err
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
		return SlotResponseDto{}, err
	}

	// Update event status
	event := model.Event{
		Id:     selectedSlot.EventId,
		Status: constants.EVENT_STATUS_UPCOMING,
	}
	if err := s.eventRepository.Updates(&event); err != nil {
		return SlotResponseDto{}, err
	}

	// Send validated slot to all participants via SSE
	s.sseService.BroadcastSlotsUpdate(slot.EventId, sse.MapToSSESlotUpdateMessages(constants.SLOT_ACTION_CREATED, []model.Slot{slot}))

	// Send event confirmation emails to all participants (including owner)
	var participants []model.Account
	if err := s.accountEventRepository.FindAccountsByEventId(event.Id, &participants); err != nil {
		log.Error().Err(err).Str("eventId", event.Id.String()).Msg("Failed to get participants for event confirmation mail")
	} else {
		for _, participant := range participants {
			go s.mailService.SendEventConfirmationEmail(
				participant,
				event,
				event.Id,
				event.OwnerId,
				slot.StartsAt,
				slot.EndsAt,
			)
		}
	}

	return MapToSlotResponseDto(slot), nil
}

func (s *SlotService) RemoveValidatedSlot(slotId uuid.UUID, userId uuid.UUID) error {
	var selectedSlot model.Slot
	if err := s.slotRepository.FindOneById(slotId, &selectedSlot); err != nil {
		return constants.ERR_SLOT_NOT_FOUND.Err
	}

	// Check if user is admin of the event
	if !selectedSlot.Event.IsOwner(&userId) {
		return constants.ERR_EVENT_ACCESS_DENIED.Err
	}

	if !selectedSlot.IsValidated {
		return constants.ERR_SLOT_NOT_FOUND.Err
	}

	// Remove validated slot
	err := s.slotRepository.DeleteValidatedSlotByEventId(selectedSlot.EventId)
	if err != nil {
		return err
	}

	// Update event status
	event := model.Event{
		Id:     selectedSlot.EventId,
		Status: constants.EVENT_STATUS_IN_DECISION,
	}
	if err := s.eventRepository.Updates(&event); err != nil {
		return err
	}

	// Send removed validated slot to all participants via SSE
	s.sseService.BroadcastSlotsUpdate(selectedSlot.EventId, sse.MapToSSESlotUpdateMessages(constants.SLOT_ACTION_DELETED, []model.Slot{selectedSlot}))

	// Send cancellation emails
	var participants []model.Account
	if err := s.accountEventRepository.FindAccountsByEventId(event.Id, &participants); err != nil {
		log.Error().Err(err).Str("eventId", event.Id.String()).Msg("Failed to get participants for event cancellation mail")
	} else {
		for _, participant := range participants {
			go s.mailService.SendEventCancellationEmail(
				participant,
				selectedSlot.Event,
				event.Id,
				selectedSlot.Event.OwnerId,
				selectedSlot.StartsAt,
				selectedSlot.EndsAt,
			)
		}
	}

	// Recalculate slots
	s.LoadSlotsAsync(selectedSlot.EventId)

	return nil
}

// LoadSlotsAsync runs LoadSlots in the background and tracks it so callers
// (tests) can wait for it deterministically with WaitAsyncLoads.
func (s *SlotService) LoadSlotsAsync(eventId uuid.UUID) {
	s.asyncLoads.Add(1)
	go func() {
		defer s.asyncLoads.Done()
		s.LoadSlots(eventId)
	}()
}

// WaitAsyncLoads blocks until every LoadSlotsAsync call has finished.
func (s *SlotService) WaitAsyncLoads() {
	s.asyncLoads.Wait()
}

// Recalculates the slots of an event, then only persists and broadcasts the ones that changed
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
	if hasStatus, err := event.CheckAndAutoUpdateStatus(s.eventRepository.Updates, &[]constants.EventStatus{constants.EVENT_STATUS_IN_DECISION}); !hasStatus || err != nil {
		if err != nil {
			log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to check event status")
		} else {
			log.Debug().Str("eventId", eventId.String()).Msg("Event is locked, skipping slot recalculation")
		}
		return
	}

	// Get existing slots for this event
	var existingSlots []model.Slot
	if err := s.slotRepository.FindByEventId(eventId, &existingSlots); err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to get existing slots")
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

	// Find common available time slots (none if less than 2 active users), existing ones that disappear get deleted
	commonSlots := s.findIntersectingTimeSlots(userAvailabilities, time.Duration(event.Duration)*time.Minute)

	// Only keep the slots that changed
	changes := s.diffSlots(eventId, existingSlots, commonSlots)
	if len(changes.Created) == 0 && len(changes.Updated) == 0 && len(changes.Deleted) == 0 {
		log.Debug().Str("eventId", eventId.String()).Msg("Slots are already up to date")
		return
	}

	if err := s.slotRepository.ApplyChanges(changes.Created, changes.Updated, changes.Deleted); err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to apply slot changes")
		return
	}

	log.Debug().
		Str("eventId", eventId.String()).
		Int("slotsCreated", len(changes.Created)).
		Int("slotsUpdated", len(changes.Updated)).
		Int("slotsDeleted", len(changes.Deleted)).
		Msg("Slot recalculation completed")

	// Send slot changes to all participants via SSE, deletions first so that applying them in order never overlaps two slots
	messages := sse.MapToSSESlotUpdateMessages(constants.SLOT_ACTION_DELETED, changes.Deleted)
	messages = append(messages, sse.MapToSSESlotUpdateMessages(constants.SLOT_ACTION_UPDATED, changes.Updated)...)
	messages = append(messages, sse.MapToSSESlotUpdateMessages(constants.SLOT_ACTION_CREATED, changes.Created)...)
	s.sseService.BroadcastSlotsUpdate(eventId, messages)
}

// Compares the computed time slots with the existing slots of an event.
// Each computed slot reuses the first existing slot it overlaps that is not reused yet (updated if its dates changed),
// or is created if there is none. Existing slots that are not reused are deleted (e.g. the second slot of a merge).
func (s *SlotService) diffSlots(eventId uuid.UUID, existingSlots []model.Slot, computedSlots []TimeSlot) SlotChanges {
	// Sort existing slots by start time (computed slots already are)
	sort.Slice(existingSlots, func(i, j int) bool {
		return existingSlots[i].StartsAt.Before(existingSlots[j].StartsAt)
	})

	var changes SlotChanges
	reused := make([]bool, len(existingSlots))
	for _, computed := range computedSlots {
		match := -1
		for i, existing := range existingSlots {
			if !reused[i] && existing.StartsAt.Before(computed.EndsAt) && computed.StartsAt.Before(existing.EndsAt) {
				match = i
				break
			}
		}

		// No existing slot to reuse, create a new one
		if match == -1 {
			changes.Created = append(changes.Created, model.Slot{
				Id:          uuid.New(),
				EventId:     eventId,
				StartsAt:    computed.StartsAt,
				EndsAt:      computed.EndsAt,
				IsValidated: false,
			})
			continue
		}

		reused[match] = true
		slot := existingSlots[match]
		if slot.StartsAt.Equal(computed.StartsAt) && slot.EndsAt.Equal(computed.EndsAt) {
			continue
		}

		slot.StartsAt = computed.StartsAt
		slot.EndsAt = computed.EndsAt
		changes.Updated = append(changes.Updated, slot)
	}

	for i, existing := range existingSlots {
		if !reused[i] {
			changes.Deleted = append(changes.Deleted, existing)
		}
	}

	return changes
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
