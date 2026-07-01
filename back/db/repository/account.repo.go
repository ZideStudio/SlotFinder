package repository

import (
	"app/commons/constants"
	"app/db"
	model "app/db/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(database *gorm.DB) *AccountRepository {
	if database == nil {
		database = db.GetDB()
	}
	return &AccountRepository{
		db: database,
	}
}

type AccountCreateDto struct {
	Id           uuid.UUID
	UserName     *string
	Email        *string
	Color        string
	Language     constants.AccountLanguage
	Password     string
	AvatarUrl    string
	AvatarData   []byte
	TermsVersion *string
	Providers    []model.AccountProvider
	TimeZone     time.Location
}

func (r *AccountRepository) Create(data AccountCreateDto, account *model.Account) error {
	*account = model.Account{
		Id:           data.Id,
		UserName:     data.UserName,
		Email:        data.Email,
		Color:        data.Color,
		Language:     data.Language,
		AvatarUrl:    data.AvatarUrl,
		AvatarData:   data.AvatarData,
		Providers:    data.Providers,
		TermsVersion: data.TermsVersion,
		TimeZone:     data.TimeZone.String(),
	}
	if account.TermsVersion != nil {
		now := time.Now().UTC()
		account.TermsAcceptedAt = &now
	}

	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::CREATE Failed to hash password")
			return err
		}
		hashedPasswordToString := string(hashedPassword)
		account.Password = &hashedPasswordToString
	}

	if err := r.db.Create(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::CREATE Failed to create account")
		return err
	}

	return nil
}

func (r *AccountRepository) Updates(account model.Account) error {
	if account.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*account.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE Failed to hash password")
			return err
		}
		hashedPasswordToString := string(hashedPassword)
		account.Password = &hashedPasswordToString
	}

	if err := r.db.Model(&account).Omit(clause.Associations).Updates(account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE Failed to update account")
		return err
	}

	return nil
}

func (r *AccountRepository) FindOneById(id uuid.UUID, account *model.Account) error {
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id.String()).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_ID Failed to find account by id")
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindOneByUsername(username string, account *model.Account, excludeId *uuid.UUID) error {
	query := r.db.Where("LOWER(username) = LOWER(?) AND deleted_at IS NULL", username)

	if excludeId != nil {
		query = query.Where("id != ?", excludeId.String())
	}

	if err := query.Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_USERNAME Failed to find account by username")
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindOneByEmail(email string, account *model.Account) error {
	if err := r.db.Where("LOWER(email) = LOWER(?) AND deleted_at IS NULL", email).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_EMAIL Failed to find account by email")
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindOneByEmailOrUsername(emailOrUsername string, account *model.Account) error {
	if err := r.db.Where("(LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)) AND deleted_at IS NULL", emailOrUsername, emailOrUsername).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_EMAIL_OR_USERNAME Failed to find account by email or username")
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindOneByResetToken(resetToken string, account *model.Account) error {
	if err := r.db.Where("reset_token = ? AND deleted_at IS NULL", resetToken).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_RESET_TOKEN Failed to find account by reset token")
		}
		return err
	}
	return nil
}

func (r *AccountRepository) UpdateResetToken(id uuid.UUID, resetToken *string, resetTokenAt *time.Time) error {
	updates := map[string]interface{}{
		"reset_token":             resetToken,
		"password_reset_token_at": resetTokenAt,
	}

	if err := r.db.Model(&model.Account{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE_RESET_TOKEN Failed to update reset token")
		return err
	}

	return nil
}

func (r *AccountRepository) FindAvatarById(id uuid.UUID) ([]byte, *time.Time, error) {
	var result struct {
		AvatarData []byte
		UpdatedAt  time.Time
	}
	if err := r.db.Model(&model.Account{}).Select("avatar_data, updated_at").Where("id = ? AND deleted_at IS NULL", id.String()).Scan(&result).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_AVATAR_BY_ID Failed to find avatar")
		return nil, nil, err
	}
	return result.AvatarData, &result.UpdatedAt, nil
}

func (r *AccountRepository) Delete(id uuid.UUID) error {
	var account model.Account
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id.String()).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::DELETE Failed to find account by id")
		}
		return err
	}

	if err := r.db.Delete(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::DELETE Failed to delete account")
		}
		return err
	}

	return nil
}
