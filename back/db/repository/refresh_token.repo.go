package repository

import (
	"app/db"
	model "app/db/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct{}

// GenerateRefreshToken generates a cryptographically secure random token
func (*RefreshTokenRepository) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::GENERATE Failed to generate random token")
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashToken creates a SHA-256 hash of the token for storage
func (*RefreshTokenRepository) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Create creates a new refresh token in the database
func (r *RefreshTokenRepository) Create(accountId uuid.UUID, expiresAt time.Time) (string, error) {
	token, err := r.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	tokenHash := r.HashToken(token)

	refreshToken := model.RefreshToken{
		Id:        uuid.New(),
		AccountId: accountId,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	if err := db.GetDB().Create(&refreshToken).Error; err != nil {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::CREATE Failed to create refresh token")
		return "", err
	}

	return token, nil
}

// FindByTokenHash finds a refresh token by its hash
func (r *RefreshTokenRepository) FindByTokenHash(tokenHash string, refreshToken *model.RefreshToken) error {
	err := db.GetDB().Where("token_hash = ? AND is_revoked = ?", tokenHash, false).First(&refreshToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::FIND_BY_TOKEN_HASH Failed to find refresh token")
		return err
	}
	return err
}

// Revoke marks a refresh token as revoked
func (*RefreshTokenRepository) Revoke(id uuid.UUID) error {
	now := time.Now()
	err := db.GetDB().Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error

	if err != nil {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::REVOKE Failed to revoke refresh token")
		return err
	}

	return nil
}

// RevokeAllForAccount revokes all refresh tokens for a specific account
func (*RefreshTokenRepository) RevokeAllForAccount(accountId uuid.UUID) error {
	now := time.Now()
	err := db.GetDB().Model(&model.RefreshToken{}).
		Where("account_id = ? AND is_revoked = ?", accountId, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error

	if err != nil {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::REVOKE_ALL_FOR_ACCOUNT Failed to revoke all refresh tokens")
		return err
	}

	return nil
}

// DeleteExpired removes expired refresh tokens from the database
func (*RefreshTokenRepository) DeleteExpired() error {
	err := db.GetDB().Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
		Delete(&model.RefreshToken{}).Error

	if err != nil {
		log.Error().Err(err).Msg("REFRESH_TOKEN_REPOSITORY::DELETE_EXPIRED Failed to delete expired refresh tokens")
		return err
	}

	return nil
}
