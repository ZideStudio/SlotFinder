package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHexa(t *testing.T) {
	t.Parallel()
	assert.True(t, IsHexa("#FFFFFF"))
	assert.True(t, IsHexa("#abc123"))
	assert.False(t, IsHexa("#FFF"))
	assert.False(t, IsHexa("FFFFFF"))
	assert.False(t, IsHexa("#GGGGGG"))
	assert.False(t, IsHexa(""))
}

func TestBoolToString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "true", BoolToString(true))
	assert.Equal(t, "false", BoolToString(false))
}

func TestCapitalize(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Hello", Capitalize("hello"))
	assert.Equal(t, "Hello", Capitalize("Hello"))
	assert.Equal(t, "", Capitalize(""))
	assert.Equal(t, "É", Capitalize("é"))
}
