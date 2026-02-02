package signin

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/config"
	model "app/db/models"
	"app/db/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type SigninService struct {
	accountRepository *repository.AccountRepository
	config            *config.Config
}

func NewSigninService(service *SigninService) *SigninService {
	if service != nil {
		return service
	}

	return &SigninService{
		accountRepository: &repository.AccountRepository{},
		config:            config.GetConfig(),
	}
}

func (s *SigninService) Signin(data *SigninDto) (token TokenResponseDto, err error) {
	var account model.Account
	if err := s.accountRepository.FindOneByEmailOrUsername(data.Identifier, &account); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return token, constants.ERR_INVALID_IDENTIFIER_OR_PASSWORD.Err
		}
		return token, err
	}

	if !account.ComparePassword(data.Password) {
		return token, constants.ERR_INVALID_IDENTIFIER_OR_PASSWORD.Err
	}

	claims := &guard.Claims{
		Id:            account.Id,
		Username:      account.UserName,
		Email:         account.Email,
		TermsAccepted: account.TermsAcceptedAt != nil,
	}

	return s.GenerateToken(claims)
}

func (s *SigninService) GenerateToken(claims *guard.Claims) (tokenResponse TokenResponseDto, err error) {
	privateKeyFile, err := os.ReadFile(s.config.Auth.PrivatePemPath)
	if err != nil {
		return tokenResponse, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return tokenResponse, err
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(168 * time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return tokenResponse, err
	}
	tokenResponse.AccessToken = tokenString

	return tokenResponse, nil
}
