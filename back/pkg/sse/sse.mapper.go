package sse

import (
	"app/commons/constants"
	model "app/db/models"
)

func MapToSSESlotUpdateMessages(action constants.SlotAction, slots []model.Slot) []SSESlotUpdateMessage {
	messages := make([]SSESlotUpdateMessage, 0, len(slots))
	for _, s := range slots {
		messages = append(messages, SSESlotUpdateMessage{
			Action:      action,
			Id:          s.Id,
			StartsAt:    s.StartsAt,
			EndsAt:      s.EndsAt,
			IsValidated: s.IsValidated,
		})
	}

	return messages
}
