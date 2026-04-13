package constants

import (
	"database/sql/driver"
	"fmt"
)

type EventStatus string

const (
	EVENT_STATUS_IN_DECISION EventStatus = "IN_DECISION"
	EVENT_STATUS_UPCOMING    EventStatus = "UPCOMING"
	EVENT_STATUS_FINISHED    EventStatus = "FINISHED"
)

var EventStatuses = []EventStatus{EVENT_STATUS_IN_DECISION, EVENT_STATUS_UPCOMING, EVENT_STATUS_FINISHED}

func (ct *EventStatus) Scan(value any) error {
	if value == nil {
		*ct = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*ct = EventStatus(string(v))
		return nil
	case string:
		*ct = EventStatus(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into EventStatus", value)
	}
}

func (ct EventStatus) Value() (driver.Value, error) {
	fmt.Println("B")
	return string(ct), nil
}
