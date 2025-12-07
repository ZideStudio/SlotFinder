package slot

import (
	"time"
)

type ConfirmSlotDto struct {
	StartsAt time.Time `json:"startsAt" binding:"required"`
	EndsAt   time.Time `json:"endsAt" binding:"required"`
}
