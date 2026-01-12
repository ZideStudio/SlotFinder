package repository

import (
	"app/commons/constants"
	"app/db"
	model "app/db/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EventRepository struct{}

func (*EventRepository) Create(event *model.Event) error {
	if err := db.GetDB().Create(&event).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::CREATE Failed to create event")
		return err
	}

	event.Sanitized()

	return nil
}

func (r *EventRepository) Updates(event *model.Event) error {
	if err := db.GetDB().Updates(&event).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::UPDATE Failed to update event")
		return err
	}

	return r.FindOneById(event.Id, event)
}

func (r *EventRepository) FindOneById(
	eventId uuid.UUID,
	event *model.Event,
) error {
	db := db.GetDB()

	if err := db.
		Where("event.id = ?", eventId).
		Preload("Owner").
		Preload("Slots").
		Preload("Availabilities").
		Preload("Availabilities.Account").
		Preload("AccountEvents.Account").
		First(&event).
		Error; err != nil {
		return err
	}

	participants := make([]model.Account, 0, len(event.AccountEvents))
	for _, ae := range event.AccountEvents {
		account := ae.Account.Sanitized(ae.Color)
		if event.OwnerId == ae.Account.Id {
			event.Owner = account
		}

		participants = append(participants, account)
	}

	event.Participants = participants

	return nil
}

func (r *EventRepository) FindEventsByAccountId(
	accountId uuid.UUID,
	limit int,
	offset int,
) ([]model.Event, int64, error) {
	db := db.GetDB()

	// Count
	var total int64
	if err := db.
		Table("account_event").
		Select("COUNT(DISTINCT event_id)").
		Where("account_id = ?", accountId).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	// Get paginated event IDs
	subQuery := db.
		Table("event").
		Select(`
			event.id,
			event.name,
			event.status,
			CASE
				WHEN event.status = ? THEN 1
				WHEN event.status = ? THEN 2
				WHEN event.status = ? THEN 3
				ELSE 4
			END AS sort_order
		`,
			constants.EVENT_STATUS_IN_DECISION,
			constants.EVENT_STATUS_UPCOMING,
			constants.EVENT_STATUS_FINISHED,
		).
		Joins("JOIN account_event ae ON ae.event_id = event.id").
		Where("ae.account_id = ?", accountId)

	var eventIds []uuid.UUID
	if err := db.
		Table("(?) AS e", subQuery).
		Order("e.sort_order ASC").
		Order("e.name ASC").
		Order("e.id ASC").
		Limit(limit).
		Offset(offset).
		Pluck("e.id", &eventIds).
		Error; err != nil {
		return nil, 0, err
	}

	if len(eventIds) == 0 {
		return []model.Event{}, total, nil
	}

	// Fetch full event details
	var events []model.Event
	if err := db.
		Where("id IN ?", eventIds).
		Preload("Owner").
		Preload("AccountEvents.Account").
		Find(&events).
		Error; err != nil {
		return nil, 0, err
	}

	// Map events by ID
	eventById := make(map[uuid.UUID]*model.Event, len(events))
	for i := range events {
		eventById[events[i].Id] = &events[i]
	}

	for _, event := range events {
		participants := make([]model.Account, 0, len(event.AccountEvents))
		for _, ae := range event.AccountEvents {
			if ae.AccountId == event.OwnerId {
				eventById[event.Id].Owner = event.Owner.Sanitized(ae.Color)
			}
			participants = append(participants, ae.Account.Sanitized(ae.Color))
		}
		eventById[event.Id].Participants = participants
	}

	orderedEvents := make([]model.Event, 0, len(eventIds))
	for _, id := range eventIds {
		if event, ok := eventById[id]; ok {
			orderedEvents = append(orderedEvents, *event)
		}
	}

	return orderedEvents, total, nil
}

func (*EventRepository) Delete(id uuid.UUID) error {
	if err := db.GetDB().Where("id = ?", id.String()).Delete(&model.Event{}).Error; err != nil {
		log.Error().Err(err).Msg("EVENT_REPOSITORY::DELETE Failed to delete event")
		return err
	}

	return nil
}
