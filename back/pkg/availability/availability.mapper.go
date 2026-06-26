package availability

import model "app/db/models"

func MapToAvailabilityResponseDto(a model.Availability) AvailabilityResponseDto {
	return AvailabilityResponseDto{
		Id:       a.Id,
		StartsAt: a.StartsAt,
		EndsAt:   a.EndsAt,
	}
}
