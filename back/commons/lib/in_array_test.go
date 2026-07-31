package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray_Found(t *testing.T) {
	assert.Equal(t, 1, InArray("b", []string{"a", "b", "c"}))
}

func TestInArray_NotFound(t *testing.T) {
	assert.Equal(t, -1, InArray("z", []string{"a", "b", "c"}))
}

func TestInArray_EmptySlice(t *testing.T) {
	assert.Equal(t, -1, InArray(1, []int{}))
}
