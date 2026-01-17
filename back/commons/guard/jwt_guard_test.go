package guard

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestShouldRenewToken(t *testing.T) {
	// Test token that should be renewed (3 minutes until expiry)
	claims1 := &Claims{}
	claims1.ExpiresAt = jwt.NewNumericDate(time.Now().Add(3 * time.Minute))
	assert.True(t, ShouldRenewToken(claims1), "Expected token with 3 minutes remaining to be renewed")

	// Test token that should not be renewed (10 minutes until expiry)
	claims2 := &Claims{}
	claims2.ExpiresAt = jwt.NewNumericDate(time.Now().Add(10 * time.Minute))
	assert.False(t, ShouldRenewToken(claims2), "Expected token with 10 minutes remaining NOT to be renewed")

	// Test expired token
	claims3 := &Claims{}
	claims3.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
	assert.False(t, ShouldRenewToken(claims3), "Expected expired token NOT to be renewed")
}
