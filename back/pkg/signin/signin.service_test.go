package signin

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/config"
	appdb "app/db"
	model "app/db/models"
	"app/db/repository"
	"app/testutils"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSigninService(t *testing.T) *SigninService {
	t.Helper()
	testutils.TestDB(t)
	return &SigninService{
		accountRepository:      repository.NewAccountRepository(nil),
		refreshTokenRepository: repository.NewRefreshTokenRepository(nil),
		config:                 config.GetConfig(),
	}
}

func createTestAccount(t *testing.T, s *SigninService, username, password string) uuid.UUID {
	t.Helper()
	id := uuid.New()
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{
		Id:       id,
		Username: &username,
		Password: password,
	}, &account))
	return id
}

func TestNewSigninService_ReusesProvidedInstance(t *testing.T) {
	existing := &SigninService{}
	assert.Same(t, existing, NewSigninService(existing))
}

func TestNewSigninService_Nil_BuildsRealDependencies(t *testing.T) {
	s := NewSigninService(nil)
	assert.NotNil(t, s.accountRepository)
	assert.NotNil(t, s.refreshTokenRepository)
}

func TestSignin_AccountRepositoryError(t *testing.T) {
	s := newTestSigninService(t)
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))

	_, err := s.Signin(&SigninDto{Identifier: "someone", Password: "whatever"})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, constants.ERR_INVALID_IDENTIFIER_OR_PASSWORD.Err)
}

func TestSignin_InvalidIdentifier(t *testing.T) {
	s := newTestSigninService(t)
	_, err := s.Signin(&SigninDto{Identifier: "unknown-" + uuid.NewString(), Password: "whatever"})
	assert.ErrorIs(t, err, constants.ERR_INVALID_IDENTIFIER_OR_PASSWORD.Err)
}

func TestSignin_WrongPassword(t *testing.T) {
	s := newTestSigninService(t)
	username := "signin-" + uuid.NewString()
	createTestAccount(t, s, username, "CorrectPassword123!")

	_, err := s.Signin(&SigninDto{Identifier: username, Password: "WrongPassword"})
	assert.ErrorIs(t, err, constants.ERR_INVALID_IDENTIFIER_OR_PASSWORD.Err)
}

func TestSignin_Success(t *testing.T) {
	s := newTestSigninService(t)
	username := "signinok-" + uuid.NewString()
	createTestAccount(t, s, username, "CorrectPassword123!")

	tokens, err := s.Signin(&SigninDto{Identifier: username, Password: "CorrectPassword123!"})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestGenerateAccessToken(t *testing.T) {
	s := newTestSigninService(t)
	token, err := s.GenerateAccessToken(&guard.Claims{Id: uuid.New()})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokens(t *testing.T) {
	s := newTestSigninService(t)
	accountId := createTestAccount(t, s, "gentokens-"+uuid.NewString(), "")
	tokens, err := s.GenerateTokens(&guard.Claims{Id: accountId})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestGenerateTokens_AccessTokenFails(t *testing.T) {
	s := newTestSigninService(t)
	cfgCopy := *s.config
	cfgCopy.Auth.PrivatePemPath = "/nonexistent/private.pem"
	s.config = &cfgCopy

	_, err := s.GenerateTokens(&guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestGenerateTokens_RefreshTokenRepositoryError(t *testing.T) {
	s := newTestSigninService(t)
	s.refreshTokenRepository = repository.NewRefreshTokenRepository(testutils.ClosedDB(t))

	_, err := s.GenerateTokens(&guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestGenerateAccessToken_PrivateKeyFileMissing(t *testing.T) {
	s := newTestSigninService(t)
	cfgCopy := *s.config
	cfgCopy.Auth.PrivatePemPath = "/nonexistent/private.pem"
	s.config = &cfgCopy

	_, err := s.GenerateAccessToken(&guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestGenerateAccessToken_PrivateKeyFileInvalid(t *testing.T) {
	s := newTestSigninService(t)
	cfgCopy := *s.config
	invalidPath := filepath.Join(t.TempDir(), "invalid.pem")
	require.NoError(t, os.WriteFile(invalidPath, []byte("not a pem file"), 0o600))
	cfgCopy.Auth.PrivatePemPath = invalidPath
	s.config = &cfgCopy

	_, err := s.GenerateAccessToken(&guard.Claims{Id: uuid.New()})
	assert.Error(t, err)
}

func TestRefreshAccessToken_RefreshTokenRepositoryError(t *testing.T) {
	s := newTestSigninService(t)
	s.refreshTokenRepository = repository.NewRefreshTokenRepository(testutils.ClosedDB(t))

	_, err := s.RefreshAccessToken("some-token")
	assert.Error(t, err)
	assert.NotEqual(t, "invalid refresh token", err.Error())
}

func TestRefreshAccessToken_InvalidToken(t *testing.T) {
	s := newTestSigninService(t)
	_, err := s.RefreshAccessToken("not-a-real-token")
	assert.Error(t, err)
}

func TestRefreshAccessToken_ExpiredToken(t *testing.T) {
	s := newTestSigninService(t)
	accountId := createTestAccount(t, s, "expiredtoken-"+uuid.NewString(), "")
	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(-time.Hour))
	require.NoError(t, err)

	_, err = s.RefreshAccessToken(rawToken)
	assert.Error(t, err)
}

func TestRefreshAccessToken_RevokedToken(t *testing.T) {
	s := newTestSigninService(t)
	accountId := createTestAccount(t, s, "revokedtoken-"+uuid.NewString(), "")
	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(time.Hour))
	require.NoError(t, err)

	var stored model.RefreshToken
	require.NoError(t, s.refreshTokenRepository.FindByTokenHash(s.refreshTokenRepository.HashToken(rawToken), &stored))
	require.NoError(t, s.refreshTokenRepository.Revoke(stored.Id))

	_, err = s.RefreshAccessToken(rawToken)
	assert.Error(t, err)
}

func TestRefreshAccessToken_Success(t *testing.T) {
	s := newTestSigninService(t)
	username := "refresh-" + uuid.NewString()
	accountId := createTestAccount(t, s, username, "")

	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(time.Hour))
	require.NoError(t, err)

	tokens, err := s.RefreshAccessToken(rawToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)

	// Token rotation: the old refresh token should now be revoked.
	// (FindByTokenHash itself filters out revoked tokens, so query directly.)
	var oldToken model.RefreshToken
	require.NoError(t, appdb.GetDB().Where("token_hash = ?", s.refreshTokenRepository.HashToken(rawToken)).First(&oldToken).Error)
	assert.True(t, oldToken.IsRevoked)
}

func TestRefreshAccessToken_RevokeFails(t *testing.T) {
	s := newTestSigninService(t)
	accountId := createTestAccount(t, s, "revokefails-"+uuid.NewString(), "")
	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(time.Hour))
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	_, err = s.RefreshAccessToken(rawToken)
	assert.Error(t, err)
}

// Forces the account-lookup branch via a closed repository: a token for a
// nonexistent account can't be created directly anymore — refresh_token.account_id
// has an ON DELETE CASCADE foreign key, so it can never outlive its account.
func TestRefreshAccessToken_AccountLookupFails(t *testing.T) {
	s := newTestSigninService(t)
	accountId := createTestAccount(t, s, "accountlookupfails-"+uuid.NewString(), "")
	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(time.Hour))
	require.NoError(t, err)

	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))

	_, err = s.RefreshAccessToken(rawToken)
	assert.Error(t, err)
}
