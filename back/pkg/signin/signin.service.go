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
	accountRepository      *repository.AccountRepository
	refreshTokenRepository *repository.RefreshTokenRepository
	config                 *config.Config
}

func NewSigninService(service *SigninService) *SigninService {
	if service != nil {
		return service
	}

	return &SigninService{
		accountRepository:      &repository.AccountRepository{},
		refreshTokenRepository: &repository.RefreshTokenRepository{},
		config:                 config.GetConfig(),
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

	return s.GenerateTokens(claims)
}

func (s *SigninService) GenerateTokens(claims *guard.Claims) (tokenResponse TokenResponseDto, err error) {
	// Generate access token
	accessToken, err := s.GenerateAccessToken(claims)
	if err != nil {
		return tokenResponse, err
	}

	// Generate refresh token
	refreshToken, err := s.refreshTokenRepository.Create(
		claims.Id,
		time.Now().Add(constants.REFRESH_TOKEN_EXPIRATION),
	)
	if err != nil {
		return tokenResponse, err
	}

	tokenResponse.AccessToken = accessToken
	tokenResponse.RefreshToken = refreshToken

	return tokenResponse, nil
}

func (s *SigninService) GenerateAccessToken(claims *guard.Claims) (string, error) {
	privateKeyFile, err := os.ReadFile(s.config.Auth.PrivatePemPath)
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return "", err
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(constants.ACCESS_TOKEN_EXPIRATION))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *SigninService) RefreshAccessToken(refreshTokenString string) (tokenResponse TokenResponseDto, err error) {
	// Hash the refresh token to look it up in the database
	tokenHash := s.refreshTokenRepository.HashToken(refreshTokenString)

	// Find the refresh token in the database
	var refreshToken model.RefreshToken
	if err := s.refreshTokenRepository.FindByTokenHash(tokenHash, &refreshToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tokenResponse, errors.New("invalid refresh token")
		}
		return tokenResponse, err
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return tokenResponse, errors.New("refresh token expired")
	}

	// Check if token is revoked
	if refreshToken.IsRevoked {
		return tokenResponse, errors.New("refresh token revoked")
	}

	// Get the account
	var account model.Account
	if err := s.accountRepository.FindOneById(refreshToken.AccountId, &account); err != nil {
		return tokenResponse, err
	}

	// Revoke the old refresh token (token rotation)
	if err := s.refreshTokenRepository.Revoke(refreshToken.Id); err != nil {
		return tokenResponse, err
	}

	// Generate new tokens
	claims := &guard.Claims{
		Id:       account.Id,
		Username: account.UserName,
		Email:    account.Email,
	}

	return s.GenerateTokens(claims)
}
