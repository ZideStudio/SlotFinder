package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	t.Parallel()
	assert.True(t, IsValidEmail("user@example.com"))
	assert.True(t, IsValidEmail("user.name+tag@example.co.uk"))
	assert.False(t, IsValidEmail("not-an-email"))
	assert.False(t, IsValidEmail("missing-domain@"))
	assert.False(t, IsValidEmail("@missing-local.com"))
	assert.False(t, IsValidEmail(""))
}
