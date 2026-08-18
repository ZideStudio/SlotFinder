package lib

import (
	"app/commons/constants"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCustomError_True(t *testing.T) {
	t.Parallel()
	assert.True(t, IsCustomError(constants.ERR_EVENT_NOT_FOUND.Err))
}

func TestIsCustomError_False(t *testing.T) {
	t.Parallel()
	assert.False(t, IsCustomError(errors.New("some random error")))
}
