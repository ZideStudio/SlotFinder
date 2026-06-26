package slot

import model "app/db/models"

func MapToSlotResponseDto(s model.Slot) SlotResponseDto {
	return SlotResponseDto{
		Id:          s.Id,
		IsValidated: s.IsValidated,
		StartsAt:    s.StartsAt,
		EndsAt:      s.EndsAt,
	}
}
