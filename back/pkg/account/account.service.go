package account

import (
	"app/commons/constants"
	"app/commons/encryption"
	"app/commons/guard"
	"app/commons/lib"
	"app/config"
	"app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
	"app/pkg/signin"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AccountService struct {
	accountRepository     *repository.AccountRepository
	avatarService         *AvatarService
	signinService         *signin.SigninService
	mailService           *mail.MailService
	config                *config.Config
	passwordResetCooldown *cache.Cache
}

func NewAccountService(service *AccountService) *AccountService {
	if service != nil {
		return service
	}

	return &AccountService{
		accountRepository:     &repository.AccountRepository{},
		avatarService:         NewAvatarService(),
		signinService:         signin.NewSigninService(nil),
		mailService:           mail.NewMailService(nil),
		config:                config.GetConfig(),
		passwordResetCooldown: cache.New(10*time.Minute, 15*time.Minute),
	}
}

func (s *AccountService) CheckUserNameAvailability(userName string, excludeUserId *uuid.UUID) (bool, error) {
	var account model.Account
	if err := s.accountRepository.FindOneByUsername(userName, &account); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}

		return false, err
	}
	if excludeUserId != nil && account.Id == *excludeUserId {
		return true, nil
	}

	return false, nil
}

func (s *AccountService) Create(data *AccountCreateDto) (string, error) {
	// Validate input
	if !lib.IsValidEmail(data.Email) {
		return "", constants.ERR_INVALID_EMAIL_FORMAT.Err
	}

	if !lib.IsValidPassword(data.Password) {
		return "", constants.ERR_INVALID_PASSWORD_FORMAT.Err
	}

	// Check if email already exists
	var existingAccount model.Account
	if err := s.accountRepository.FindOneByEmail(data.Email, &existingAccount); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if existingAccount.Id != uuid.Nil {
		return "", constants.ERR_EMAIL_ALREADY_EXISTS.Err
	}

	// Choose a random color
	colors := constants.COLORS
	color := colors[mathrand.Intn(len(colors))]

	// Create account
	var account model.Account
	if err := s.accountRepository.Create(repository.AccountCreateDto{
		Email:    &data.Email,
		Color:    string(color),
		Password: data.Password,
	}, &account); err != nil {
		return "", err
	}

	// Update avatar
	if err := s.accountRepository.Updates(model.Account{
		Id:        account.Id,
		AvatarUrl: s.avatarService.GetGravatarURL(account.Id.String()),
	}); err != nil {
		return "", err
	}

	// Generate token
	claims := &guard.Claims{
		Id:       account.Id,
		Username: account.UserName,
		Email:    account.Email,
	}

	token, err := s.signinService.GenerateToken(claims)
	if err != nil {
		_ = s.accountRepository.Delete(account.Id)
		return "", err
	}

	if account.Email != nil {
		go func() {
			_ = s.mailService.SendMail(mail.EmailParams{
				Template: constants.MAIL_TEMPLATE_WELCOME,
				To:       *account.Email,
				Subject:  "Welcome to SlotFinder!",
				Params: map[string]string{
					"LoginUrl": fmt.Sprintf("%s/login", s.config.Origin),
				},
			})
		}()
	}

	return token.AccessToken, nil
}

func (s *AccountService) GetMe(userId uuid.UUID) (account model.Account, err error) {
	if err := s.accountRepository.FindOneById(userId, &account); err != nil {
		return account, err
	}

	return account, nil
}

func (s *AccountService) Update(dto *AccountUpdateDto, userId uuid.UUID) (account model.Account, accessToken *string, err error) {
	if err = s.accountRepository.FindOneById(userId, &account); err != nil {
		return account, nil, err
	}

	if dto.UserName != nil && (account.UserName == nil || *dto.UserName != *account.UserName) {
		isUserNameAvailable, err := s.CheckUserNameAvailability(*dto.UserName, &userId)
		if err != nil {
			return account, nil, err
		}
		if !isUserNameAvailable {
			return account, nil, constants.ERR_USERNAME_ALREADY_TAKEN.Err
		}
		account.UserName = dto.UserName
	}
	if dto.Email != nil {
		account.Email = dto.Email
	}
	if dto.Password != nil {
		if !lib.IsValidPassword(*dto.Password) {
			return account, nil, constants.ERR_INVALID_PASSWORD_FORMAT.Err
		}
		account.Password = dto.Password
	}
	if dto.Color != nil {
		if !lib.IsHexa(*dto.Color) {
			return account, nil, constants.ERR_INVALID_COLOR_FORMAT.Err
		}
		account.Color = *dto.Color
	}

	if err := s.accountRepository.Updates(account); err != nil {
		return account, nil, err
	}

	if dto.UserName != nil {
		claims := &guard.Claims{
			Id:       account.Id,
			Username: account.UserName,
			Email:    account.Email,
		}

		token, err := s.signinService.GenerateToken(claims)
		if err != nil {
			_ = s.accountRepository.Delete(account.Id)
			return account, nil, err
		}
		accessToken = &token.AccessToken
	}

	me, err := s.GetMe(userId)
	if err != nil {
		return account, nil, err
	}

	return me, accessToken, nil
}

// ForgotPassword generates a reset token and sends reset email
func (s *AccountService) ForgotPassword(dto *ForgotPasswordDto) error {
	// Atomic check-and-set to prevent race conditions
	if _, found := s.passwordResetCooldown.Get(dto.Email); found {
		return constants.ERR_PASSWORD_RESET_TOO_FREQUENT.Err
	}

	var account model.Account
	if err := s.accountRepository.FindOneByEmail(dto.Email, &account); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Silently ignore
		}
		return err
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := cryptorand.Read(tokenBytes); err != nil {
		return err
	}
	resetToken := hex.EncodeToString(tokenBytes)
	resetTokenEncrypted, err := encryption.Encrypt(resetToken)
	if err != nil {
		return err
	}
	expirationTime := time.Now().Add(time.Hour)

	// Update account with reset token
	if err := s.accountRepository.UpdateResetToken(account.Id, &resetToken, &expirationTime); err != nil {
		return err
	}

	// Send reset email
	if err := s.mailService.SendMail(mail.EmailParams{
		Template: constants.MAIL_TEMPLATE_PASSWORD_RESET,
		To:       *account.Email,
		Subject:  "Reset your password",
		Params: map[string]string{
			"ResetUrl":   fmt.Sprintf("%s/reset-password?token=%s", s.config.Origin, resetTokenEncrypted),
			"ExpiryTime": "1 hour",
		},
	}); err != nil {
		log.Error().Err(err).Str("email", *account.Email).Msg("ACCOUNT_SERVICE::SEND_PASSWORD_RESET_EMAIL Failed to send password reset email")
		return err
	}

	// Record this attempt
	s.passwordResetCooldown.Set(dto.Email, time.Now(), cache.DefaultExpiration)

	return nil
}

// ResetPassword validates reset token and updates password
func (s *AccountService) ResetPassword(dto *ResetPasswordDto) error {
	// Validate password format
	if !lib.IsValidPassword(dto.Password) {
		return constants.ERR_INVALID_PASSWORD_FORMAT.Err
	}

	resetTokenDecrypted, err := encryption.Decrypt(dto.Token)
	if err != nil {
		return err
	}
	fmt.Println("Reset token encrypted:", resetTokenDecrypted)

	var account model.Account
	if err := s.accountRepository.FindOneByResetToken(resetTokenDecrypted, &account); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ERR_INVALID_RESET_TOKEN.Err
		}
		return err
	}

	// Check if token has expired
	if account.PasswordResetTokenAt == nil || time.Now().After(*account.PasswordResetTokenAt) {
		return constants.ERR_RESET_TOKEN_EXPIRED.Err
	}

	// Update password and clear reset token
	account.Password = &dto.Password
	if err := s.accountRepository.Updates(account); err != nil {
		return err
	}

	// Clear reset token
	if err := s.accountRepository.UpdateResetToken(account.Id, nil, nil); err != nil {
		return err
	}

	if account.Email == nil {
		return nil
	}

	// Send confirmation email
	go func() {
		_ = s.mailService.SendMail(mail.EmailParams{
			Template: constants.MAIL_TEMPLATE_PASSWORD_RESET_CONFIRMATION,
			To:       *account.Email,
			Subject:  "Password reset successful",
			Params: map[string]string{
				"Timestamp": time.Now().Format("January 2, 2006 at 15:04 UTC"),
				"LoginUrl":  fmt.Sprintf("%s/login", s.config.Origin),
			},
		})
	}()

	return nil
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
