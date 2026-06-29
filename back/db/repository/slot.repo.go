package repository

import (
	"app/db"
	model "app/db/models"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SlotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(database *gorm.DB) *SlotRepository {
	if database == nil {
		database = db.GetDB()
	}
	return &SlotRepository{
		db: database,
	}
}

func (r *SlotRepository) Create(slot *model.Slot) error {
	if err := r.db.Create(&slot).First(&slot).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::CREATE Failed to create slot")
		return err
	}

	return nil
}

func (r *SlotRepository) Updates(slot *model.Slot) error {
	if err := r.db.Omit(clause.Associations).Updates(&slot).Error; err != nil {
		log.Error().Err(err).Msg("SLOT_REPOSITORY::UPDATES Failed to update slot")
		return err
	}

	return nil
}

func (r *SlotRepository) FindOneById(slotId uuid.UUID, slot *model.Slot) error {
	if err := r.db.Where("id = ?", slotId.String()).Preload("Event").Preload("Event.Owner").Preload("Event.AccountEvents").Preload("Event.AccountEvents.Account").First(slot).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Str("slotId", slotId.String()).Msg("SLOT_REPOSITORY::FIND_ONE_BY_ID Failed to find slot by id")
		}
		return err
	}

	return nil
}

func (r *SlotRepository) FindByEventId(eventId uuid.UUID, slots *[]model.Slot) error {
	if err := r.db.Where("event_id = ?", eventId).Find(slots).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::FIND_BY_EVENT_ID Failed to find slots by event id")
		return err
	}

	return nil
}

func (r *SlotRepository) FindValidatedSlotByEventId(eventId uuid.UUID, slot *model.Slot) error {
	if err := r.db.Where("event_id = ? AND is_validated = ?", eventId, true).Preload("Event").Preload("Event.Owner").Preload("Event.AccountEvents").Preload("Event.AccountEvents.Account").Find(slot).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::FIND_VALIDATED_SLOT_BY_EVENT_ID Failed to find validated slots by event id")
		return err
	}

	return nil
}

func (r *SlotRepository) DeleteByEventId(eventId uuid.UUID) error {
	if err := r.db.Where("event_id = ?", eventId).Delete(&model.Slot{}).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::DELETE_BY_EVENT_ID Failed to delete slots by event id")
		return err
	}

	return nil
}

func (r *SlotRepository) DeleteValidatedSlotByEventId(eventId uuid.UUID) error {
	if err := r.db.Where("event_id = ? AND is_validated = ?", eventId, true).Delete(&model.Slot{}).Error; err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("SLOT_REPOSITORY::DELETE_VALIDATED_BY_EVENT_ID Failed to delete validated slot by event ID")
		return err
	}

	return nil
}
