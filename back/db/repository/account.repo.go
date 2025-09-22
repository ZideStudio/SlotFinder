package repository

import (
	"app/db"
	model "app/db/models"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountRepository struct{}

type AccountCreateDto struct {
	UserName  string
	Email     string
	Password  string
	Providers []model.AccountProvider
}

func (*AccountRepository) Create(data AccountCreateDto, account *model.Account) error {
	*account = model.Account{
		Id:        uuid.New(),
		UserName:  data.UserName,
		Email:     data.Email,
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
	err := db.GetDB().Where("id = ? AND deleted_at IS NULL", id.String()).Preload("Providers").First(&account).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_ID Failed to find account by id")
		return err
	}
	return err
}

func (*AccountRepository) FindOneByUsername(username string, account *model.Account) error {
	err := db.GetDB().Where("LOWER(username) = LOWER(?) AND deleted_at IS NULL", username).Preload("Providers").First(&account).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_USERNAME Failed to find account by username")
		return err
	}
	return err
}

func (*AccountRepository) FindOneByEmail(email string, account *model.Account) error {
	err := db.GetDB().Where("LOWER(email) = LOWER(?) AND deleted_at IS NULL", email).Preload("Providers").First(&account).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("ACCOUNT_REPOSITORY::FIND_ONE_BY_EMAIL Failed to find account by email")
		return err
	}
	return err
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
