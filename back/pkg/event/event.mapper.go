package event

import (
	model "app/db/models"
)

// durationToFields converts total minutes to days, hours, minutes
func durationToFields(duration int) EventDurationFields {
	days := duration / (24 * 60)
	remaining := duration % (24 * 60)
	hours := remaining / 60
	minutes := remaining % 60
	return EventDurationFields{Days: days, Hours: hours, Minutes: minutes}
}

// mapToOwnerDto maps an Account to EventOwnerDto, with optional color override
func mapToOwnerDto(account model.Account, colorOverride *string) EventOwnerDto {
	color := account.Color
	if colorOverride != nil {
		color = *colorOverride
	}
	return EventOwnerDto{
		UserName:  account.UserName,
		AvatarUrl: account.AvatarUrl,
		Color:     color,
	}
}

// mapToParticipantDto maps an AccountEvent to EventParticipantDto, falls back to Account.Color
func mapToParticipantDto(ae model.AccountEvent) EventParticipantDto {
	color := ae.Account.Color
	if ae.Color != nil && *ae.Color != "" {
		color = *ae.Color
	}
	return EventParticipantDto{
		UserName:  ae.Account.UserName,
		AvatarUrl: ae.Account.AvatarUrl,
		Color:     color,
	}
}

// MapToEventListItemDto maps a model.Event to EventListItemDto
func MapToEventListItemDto(e model.Event) EventListItemDto {
	return EventListItemDto{
		Id:                  e.Id,
		Name:                e.Name,
		Description:         e.Description,
		EventDurationFields: durationToFields(e.Duration),
		StartsAt:            e.StartsAt,
		EndsAt:              e.EndsAt,
		Status:              e.Status,
	}
}

// MapToEventCreateResponseDto maps a model.Event to EventCreateResponseDto
func MapToEventCreateResponseDto(e model.Event) EventCreateResponseDto {
	return EventCreateResponseDto{
		Id:                  e.Id,
		Name:                e.Name,
		Description:         e.Description,
		EventDurationFields: durationToFields(e.Duration),
		StartsAt:            e.StartsAt,
		EndsAt:              e.EndsAt,
		Status:              e.Status,
		Owner:               mapToOwnerDto(e.Owner, nil),
	}
}

// MapToEventBasicResponseDto maps a model.Event to EventBasicResponseDto
func MapToEventBasicResponseDto(e model.Event) EventBasicResponseDto {
	return EventBasicResponseDto{
		Id:                  e.Id,
		Name:                e.Name,
		Description:         e.Description,
		EventDurationFields: durationToFields(e.Duration),
		StartsAt:            e.StartsAt,
		EndsAt:              e.EndsAt,
		Status:              e.Status,
	}
}

// MapToEventFullResponseDto maps a model.Event to EventFullResponseDto with all relations
func MapToEventFullResponseDto(e model.Event) EventFullResponseDto {
	participants := make([]EventParticipantDto, 0, len(e.AccountEvents))
	for _, ae := range e.AccountEvents {
		participants = append(participants, mapToParticipantDto(ae))
	}

	availabilities := e.Availabilities
	if availabilities == nil {
		availabilities = []model.Availability{}
	}
	for i := range availabilities {
		availabilities[i].Sanitized()
	}
	slots := e.Slots
	if slots == nil {
		slots = []model.Slot{}
	}

	return EventFullResponseDto{
		Id:                  e.Id,
		Name:                e.Name,
		Description:         e.Description,
		EventDurationFields: durationToFields(e.Duration),
		StartsAt:            e.StartsAt,
		EndsAt:              e.EndsAt,
		Status:              e.Status,
		Owner:               mapToOwnerDto(e.Owner, nil),
		Participants:        participants,
		Availabilities:      availabilities,
		Slots:               slots,
	}
}
