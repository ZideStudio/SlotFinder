package event

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- durationToFields ---

func TestDurationToFields(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		want     EventDurationFields
	}{
		{"zero", 0, EventDurationFields{Days: 0, Hours: 0, Minutes: 0}},
		{"minutes only", 45, EventDurationFields{Days: 0, Hours: 0, Minutes: 45}},
		{"hours only", 120, EventDurationFields{Days: 0, Hours: 2, Minutes: 0}},
		{"days only", 1440, EventDurationFields{Days: 1, Hours: 0, Minutes: 0}},
		{"mixed", 1500, EventDurationFields{Days: 1, Hours: 1, Minutes: 0}},
		{"max 3 weeks", 30240, EventDurationFields{Days: 21, Hours: 0, Minutes: 0}},
		{"complex", 2521, EventDurationFields{Days: 1, Hours: 18, Minutes: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, durationToFields(tt.duration))
		})
	}
}

// --- FieldsToDuration ---

func TestFieldsToDuration(t *testing.T) {
	tests := []struct {
		name    string
		days    int
		hours   int
		minutes int
		want    int
	}{
		{"zero", 0, 0, 0, 0},
		{"minutes only", 0, 0, 45, 45},
		{"hours only", 0, 2, 0, 120},
		{"days only", 1, 0, 0, 1440},
		{"mixed", 1, 1, 0, 1500},
		{"max 3 weeks", 21, 0, 0, 30240},
		{"roundtrip complex", 1, 18, 1, 2521},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FieldsToDuration(tt.days, tt.hours, tt.minutes))
		})
	}
}

// --- mapToParticipantDto ---

func TestMapToParticipantDto_UsesAccountColorByDefault(t *testing.T) {
	ae := model.AccountEvent{
		Color: nil,
		Account: model.Account{
			Color: "#AABBCC",
		},
	}
	dto := mapToParticipantDto(ae)
	assert.Equal(t, "#AABBCC", dto.Color)
}

func TestMapToParticipantDto_UsesAccountEventColorOverride(t *testing.T) {
	override := "#112233"
	ae := model.AccountEvent{
		Color: &override,
		Account: model.Account{
			Color: "#AABBCC",
		},
	}
	dto := mapToParticipantDto(ae)
	assert.Equal(t, "#112233", dto.Color)
}

// --- MapToEventListItemDto ---

func TestMapToEventListItemDto(t *testing.T) {
	now := time.Now()
	e := model.Event{
		Id:       uuid.New(),
		Name:     "My Event",
		Duration: 1500, // 1 day 1 hour
		StartsAt: now,
		EndsAt:   now.Add(48 * time.Hour),
		Status:   constants.EVENT_STATUS_IN_DECISION,
	}
	dto := MapToEventListItemDto(e)

	assert.Equal(t, e.Id, dto.Id)
	assert.Equal(t, "My Event", dto.Name)
	assert.Equal(t, 1, dto.Days)
	assert.Equal(t, 1, dto.Hours)
	assert.Equal(t, 0, dto.Minutes)
	assert.Equal(t, constants.EVENT_STATUS_IN_DECISION, dto.Status)
}

// --- MapToEventFullResponseDto ---

func TestMapToEventFullResponseDto_EmptyParticipantsWhenNoAccountEvents(t *testing.T) {
	e := model.Event{
		Id:            uuid.New(),
		Duration:      60,
		AccountEvents: nil,
		Availabilities: nil,
		Slots:         nil,
	}
	dto := MapToEventFullResponseDto(e)

	assert.NotNil(t, dto.Participants)
	assert.Len(t, dto.Participants, 0)
	assert.NotNil(t, dto.Availabilities)
	assert.NotNil(t, dto.Slots)
}

func TestMapToEventFullResponseDto_ParticipantsWithColors(t *testing.T) {
	eventColor := "#FF0000"
	userName := "alice"
	e := model.Event{
		Id:       uuid.New(),
		Duration: 60,
		AccountEvents: []model.AccountEvent{
			{
				Color: &eventColor,
				Account: model.Account{
					UserName: &userName,
					Color:    "#000000",
				},
			},
			{
				Color: nil, // should fall back to account color
				Account: model.Account{
					Color: "#FFFFFF",
				},
			},
		},
	}
	dto := MapToEventFullResponseDto(e)

	assert.Len(t, dto.Participants, 2)
	assert.Equal(t, "#FF0000", dto.Participants[0].Color)
	assert.Equal(t, "#FFFFFF", dto.Participants[1].Color)
}
