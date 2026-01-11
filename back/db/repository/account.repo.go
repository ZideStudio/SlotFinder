package repository

import (
	"app/db"
	model "app/db/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountRepository struct{}

type AccountCreateDto struct {
	UserName  *string
	Email     *string
	Color     string
	Password  string
	AvatarUrl string
	Providers []model.AccountProvider
}

func (*AccountRepository) Create(data AccountCreateDto, account *model.Account) error {
	*account = model.Account{
		Id:        uuid.New(),
		UserName:  data.UserName,
		Email:     data.Email,
		Color:     data.Color,
		AvatarUrl: data.AvatarUrl,
		Providers: data.Providers,
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

	if err := db.GetDB().Create(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::CREATE Failed to create account")
		return err
	}

	return nil
}

func (*AccountRepository) Updates(account model.Account) error {
	if account.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*account.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE Failed to hash password")
			return err
		}
		hashedPasswordToString := string(hashedPassword)
		account.Password = &hashedPasswordToString
	}

	if err := db.GetDB().Model(&account).Updates(account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE Failed to update account")
		return err
	}

	return nil
}

func (*AccountRepository) FindOneById(id uuid.UUID, account *model.Account) error {
	if err := db.GetDB().Where("id = ? AND deleted_at IS NULL", id.String()).Preload("Providers").First(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_ID Failed to find account by id")
		return err
	}
	return nil
}

func (*AccountRepository) FindOneByUsername(username string, account *model.Account, excludeId *uuid.UUID) error {
	query := db.GetDB().Where("LOWER(username) = LOWER(?) AND deleted_at IS NULL", username)

	if excludeId != nil {
		query = query.Where("id != ?", excludeId.String())
	}

	if err := query.Preload("Providers").First(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_USERNAME Failed to find account by username")
		return err
	}
	return nil
}

func (*AccountRepository) FindOneByEmail(email string, account *model.Account) error {
	if err := db.GetDB().Where("LOWER(email) = LOWER(?) AND deleted_at IS NULL", email).Preload("Providers").First(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_EMAIL Failed to find account by email")
		return err
	}
	return nil
}

func (*AccountRepository) FindOneByEmailOrUsername(emailOrUsername string, account *model.Account) error {
	if err := db.GetDB().Where("(LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)) AND deleted_at IS NULL", emailOrUsername, emailOrUsername).Preload("Providers").First(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_EMAIL_OR_USERNAME Failed to find account by email or username")
		return err
	}
	return nil
}

func (*AccountRepository) FindOneByResetToken(resetToken string, account *model.Account) error {
	if err := db.GetDB().Where("reset_token = ? AND deleted_at IS NULL", resetToken).Preload("Providers").First(&account).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_RESET_TOKEN Failed to find account by reset token")
		return err
	}
	return nil
}

func (*AccountRepository) UpdateResetToken(id uuid.UUID, resetToken *string, resetTokenAt *time.Time) error {
	updates := map[string]interface{}{
		"reset_token":             resetToken,
		"password_reset_token_at": resetTokenAt,
	}

	if err := db.GetDB().Model(&model.Account{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::UPDATE_RESET_TOKEN Failed to update reset token")
		return err
	}

	return nil
}

func (*AccountRepository) Delete(id uuid.UUID) error {
	var account model.Account
	if err := db.GetDB().Where("id = ? AND deleted_at IS NULL", id.String()).Preload("Providers").First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::DELETE Failed to find account by id")
		}
		return err
	}

	if err := db.GetDB().Delete(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::DELETE Failed to delete account")
		}
		return err
	}

	return nil
}
