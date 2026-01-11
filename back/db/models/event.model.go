package model

import (
	"app/commons/constants"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id          uuid.UUID             `gorm:"column:id;type:uuid;unique;primary_key" json:"id"`
	Name        string                `gorm:"column:name;size:255" json:"name"`
	Description *string               `gorm:"column:description;type:text" json:"description"`
	Duration    int                   `gorm:"column:duration;default:60" json:"duration"` // In minutes
	StartsAt    time.Time             `gorm:"column:starts_at" json:"startsAt"`
	EndsAt      time.Time             `gorm:"column:ends_at" json:"endsAt"`
	CreatedAt   time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	OwnerId     uuid.UUID             `gorm:"column:owner_id;type:uuid;primaryKey" json:"-"`
	Status      constants.EventStatus `gorm:"column:status" json:"status"`

	// Relations
	Owner          Account        `gorm:"foreignKey:OwnerId;references:Id" json:"owner"`
	AccountEvents  []AccountEvent `gorm:"foreignKey:EventId;references:Id" json:"-"`
	Availabilities []Availability `gorm:"foreignKey:EventId;references:Id" json:"availabilities"`
	Slots          []Slot         `gorm:"foreignKey:EventId;references:Id" json:"slots"`

	// Computed field not stored in DB
	Participants []Account `gorm:"-" json:"participants"`
}

func (Event) TableName() string {
	return "event"
}

func (e *Event) Sanitized() *Event {
	e.Owner = e.Owner.Sanitized(nil)

	if len(e.Availabilities) > 0 {
		availabilities := make([]Availability, len(e.Availabilities))
		for i, availability := range e.Availabilities {
			availabilities[i] = *availability.Sanitized()
		}

		e.Availabilities = availabilities
	}

	return e
}

// checks if a user has access to this event by verifying
// if the user ID exists in the event's AccountEvents relation.
func (e *Event) HasUserAccess(userId *uuid.UUID) bool {
	if userId == nil {
		return false
	}

	for _, accountEvent := range e.AccountEvents {
		if accountEvent.AccountId == *userId {
			return true
		}
	}

	return false
}

// checks if the given user ID has administrative privileges
func (e *Event) IsOwner(userId *uuid.UUID) bool {
	if userId == nil {
		return false
	}

	return e.OwnerId == *userId
}

// GetValidatedSlot returns the validated slot for the event
func (e *Event) GetValidatedSlot() *Slot {
	if len(e.Slots) == 0 {
		return nil
	}

	for _, slot := range e.Slots {
		if slot.IsValidated {
			return &slot
		}
	}

	return nil
}

// HasOneOfStatuses checks if the event status is one of the required statuses
func (e *Event) HasOneOfStatuses(requireOneOfStatuses *[]constants.EventStatus) bool {
	if requireOneOfStatuses == nil {
		return false
	}

	return slices.Contains(*requireOneOfStatuses, e.Status)
}

// CheckAndAutoUpdateStatus checks if the event is in decision status, and updates the status to finished if needed
// If requireOneOfStatuses is provided, it will return whether the event status is one of the required statuses
func (e *Event) CheckAndAutoUpdateStatus(updateFunc func(*Event) error, requireOneOfStatuses *[]constants.EventStatus) (hasStatus bool, err error) {
	slot := e.GetValidatedSlot()

	// Event is still in decision
	now := time.Now()
	isEventPassed := now.After(e.EndsAt)
	isValidatedSlotPassed := slot != nil && now.After(slot.EndsAt)
	if e.Status != constants.EVENT_STATUS_FINISHED && !isEventPassed && !isValidatedSlotPassed {
		return e.HasOneOfStatuses(requireOneOfStatuses), nil
	}

	// Update event status to finished
	if e.Status != constants.EVENT_STATUS_FINISHED {
		e.Status = constants.EVENT_STATUS_FINISHED
		if err := updateFunc(e); err != nil {
			return false, err
		}
	}

	return e.HasOneOfStatuses(requireOneOfStatuses), nil
}
