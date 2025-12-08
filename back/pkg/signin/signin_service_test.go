package signin

import (
	"app/db/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenRepository_GenerateRefreshToken(t *testing.T) {
	repo := &repository.RefreshTokenRepository{}
	
	token, err := repo.GenerateRefreshToken()
	assert.NoError(t, err, "Expected no error when generating refresh token")
	assert.NotEmpty(t, token, "Expected refresh token to be generated")
	assert.Greater(t, len(token), 40, "Expected token to be sufficiently long")
}

func TestRefreshTokenRepository_HashToken(t *testing.T) {
	repo := &repository.RefreshTokenRepository{}
	
	token := "test-token-12345"
	hash1 := repo.HashToken(token)
	hash2 := repo.HashToken(token)
	
	assert.NotEmpty(t, hash1, "Expected hash to be generated")
	assert.Equal(t, hash1, hash2, "Expected same token to produce same hash")
	
	differentToken := "different-token"
	hash3 := repo.HashToken(differentToken)
	assert.NotEqual(t, hash1, hash3, "Expected different tokens to produce different hashes")
}
