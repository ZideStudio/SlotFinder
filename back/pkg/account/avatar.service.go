package account

import (
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AvatarService struct {
	accountRepository *repository.AccountRepository
}

func NewAvatarService(service *AvatarService) *AvatarService {
	if service != nil {
		return service
	}

	return &AvatarService{
		accountRepository: repository.NewAccountRepository(nil),
	}
}

func GetGravatarURL(username string) string {
	hash := sha256.Sum256([]byte(username))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=retro&s=256", hashStr)
}

// FetchAndStoreGravatar fetches, processes and returns the Gravatar image bytes and local URL for the account.
// Falls back to the external Gravatar URL if the image cannot be fetched.
func (*AvatarService) FetchAndStoreGravatar(username string, accountId uuid.UUID) ([]byte, string) {
	data, err := lib.ProcessAvatarFromURL(GetGravatarURL(username))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch Gravatar image, falling back to external URL")
		return nil, GetGravatarURL(accountId.String())
	}
	return data, fmt.Sprintf("/api/v1/account/%s/avatar", accountId.String())
}

func (s *AvatarService) FindAvatarById(id uuid.UUID) ([]byte, *time.Time, error) {
	return s.accountRepository.FindAvatarById(id)
}

// UploadAvatar processes an image from a URL or raw bytes and returns the result as JPEG bytes.
func (*AvatarService) UploadAvatar(imgUrl *string, imgBytes []byte) ([]byte, error) {
	if imgUrl == nil && imgBytes == nil {
		return nil, fmt.Errorf("no image provided")
	}

	if imgUrl != nil {
		data, err := lib.ProcessAvatarFromURL(*imgUrl)
		if err != nil {
			log.Error().Err(err).Msg("Failed to process avatar from URL")
			return nil, err
		}
		return data, nil
	}

	data, err := lib.ProcessAvatar(imgBytes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process avatar bytes")
		return nil, err
	}
	return data, nil
}

func (s *AvatarService) UploadUserAvatar(imgBytes []byte, userId uuid.UUID) error {
	processed, err := s.UploadAvatar(nil, imgBytes)
	if err != nil {
		return fmt.Errorf("error processing avatar: %w", err)
	}

	avatarUrl := fmt.Sprintf("/api/v1/account/%s/avatar", userId.String())

	if err := s.accountRepository.Updates(model.Account{
		Id:         userId,
		AvatarUrl:  avatarUrl,
		AvatarData: processed,
	}); err != nil {
		return fmt.Errorf("error updating avatar on account: %w", err)
	}

	return nil
}
