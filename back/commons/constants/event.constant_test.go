package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventStatus_Scan_Nil(t *testing.T) {
	var status EventStatus = "UPCOMING"
	assert.NoError(t, status.Scan(nil))
	assert.Equal(t, EventStatus(""), status)
}

func TestEventStatus_Scan_Bytes(t *testing.T) {
	var status EventStatus
	assert.NoError(t, status.Scan([]byte("UPCOMING")))
	assert.Equal(t, EVENT_STATUS_UPCOMING, status)
}

func TestEventStatus_Scan_String(t *testing.T) {
	var status EventStatus
	assert.NoError(t, status.Scan("FINISHED"))
	assert.Equal(t, EVENT_STATUS_FINISHED, status)
}

func TestEventStatus_Scan_UnsupportedType(t *testing.T) {
	var status EventStatus
	err := status.Scan(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot scan")
}

func TestEventStatus_Value(t *testing.T) {
	status := EVENT_STATUS_IN_DECISION
	value, err := status.Value()
	assert.NoError(t, err)
	assert.Equal(t, "IN_DECISION", value)
}
