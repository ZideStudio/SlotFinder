package account

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/lib"
	"app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/signin"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
	signinService     *signin.SigninService
}

func NewAccountService(service *AccountService) *AccountService {
	if service != nil {
		return service
	}

	return &AccountService{
		accountRepository: &repository.AccountRepository{},
		signinService:     signin.NewSigninService(nil),
	}
}

func (s *AccountService) CheckUserNameAvailability(userName string) (bool, error) {
	var account model.Account
	if err := s.accountRepository.FindOneByUsername(userName, &account); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func (s *AccountService) Create(data *AccountCreateDto) (string, error) {
	if !lib.IsValidEmail(data.Email) {
		return "", constants.ERR_INVALID_EMAIL_FORMAT.Err
	}

	isUserNameAvailable, err := s.CheckUserNameAvailability(data.UserName)
	if err != nil {
		return "", err
	}
	if !isUserNameAvailable {
		return "", constants.ERR_USERNAME_ALREADY_TAKEN.Err
	}

	var existingAccount model.Account
	if err := s.accountRepository.FindOneByEmail(data.Email, &existingAccount); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if existingAccount.Id != uuid.Nil {
		return "", constants.ERR_EMAIL_ALREADY_EXISTS.Err
	}

	var account model.Account
	if err := s.accountRepository.Create(repository.AccountCreateDto{
		UserName: data.UserName,
		Email:    data.Email,
		Password: data.Password,
	}, &account); err != nil {
		return "", err
	}

	claims := &guard.Claims{
		Id:       account.Id,
		Username: account.UserName,
		Email:    account.Email,
	}

	token, err := s.signinService.GenerateToken(claims)
	if err != nil {
		s.accountRepository.Delete(account.Id)
		return "", err
	}

	return token.AccessToken, nil
}

func (s *AccountService) GetMe(userId uuid.UUID) (account model.Account, err error) {
	if err := s.accountRepository.FindOneById(userId, &account); err != nil {
		return account, err
	}

	return account, nil
}

func (s *AccountService) Update(dto *AccountUpdateDto, userId uuid.UUID) (account model.Account, err error) {
	if err = s.accountRepository.FindOneById(userId, &account); err != nil {
		return account, err
	}

	if dto.UserName != nil {
		account.UserName = *dto.UserName
	}
	if dto.Email != nil {
		account.Email = *dto.Email
	}
	if dto.Password != nil {
		account.Password = dto.Password
	}

	if err := s.accountRepository.Updates(account); err != nil {
		return account, err
	}

	return s.GetMe(userId)
}

func (s *AccountService) Delete(userId uuid.UUID, user *guard.Claims) (account model.Account, err error) {
	if err = s.accountRepository.FindOneById(userId, &account); err != nil {
		return account, err
	}

	if err = db.GetDB().Delete(&account).Error; err != nil {
		return account, err
	}

	return account, err
}
