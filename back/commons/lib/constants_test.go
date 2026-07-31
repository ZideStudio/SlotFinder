package lib

import (
	"app/commons/constants"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCustomError_True(t *testing.T) {
	assert.True(t, IsCustomError(constants.ERR_EVENT_NOT_FOUND.Err))
}

func TestIsCustomError_False(t *testing.T) {
	assert.False(t, IsCustomError(errors.New("some random error")))
}
