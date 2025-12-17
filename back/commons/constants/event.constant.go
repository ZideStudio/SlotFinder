package constants

type EventStatus string

const (
	EVENT_STATUS_IN_DECISION EventStatus = "IN_DECISION"
	EVENT_STATUS_UPCOMING    EventStatus = "UPCOMING"
	EVENT_STATUS_FINISHED    EventStatus = "FINISHED"
)
