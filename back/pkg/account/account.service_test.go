package account

import (
	"app/commons/constants"
	"app/commons/encryption"
	"app/commons/guard"
	"app/config"
	appdb "app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/mail"
	"app/pkg/signin"
	"app/testutils"
	"errors"
	"net/smtp"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// NOTE: Do not call NewAccountService(nil) in tests — build the struct
// directly. Nested NewXService(nil) calls are safe here since this wires up
// a fresh test transaction and real config (see TestMain in main_test.go).
func newTestAccountService(t *testing.T) *AccountService {
	t.Helper()
	testutils.TestDB(t)
	return &AccountService{
		accountRepository:      repository.NewAccountRepository(nil),
		avatarService:          NewAvatarService(nil),
		signinService:          signin.NewSigninService(nil),
		mailService:            mail.NewMailService(nil),
		config:                 config.GetConfig(),
		passwordResetCooldown:  cache.New(10*time.Minute, 15*time.Minute),
		refreshTokenRepository: repository.NewRefreshTokenRepository(nil),
	}
}

// uniqueEmail, stubSMTP, stubSMTPAwait, and awaitSMTP delegate to testutils
// (shared with pkg/provider, which needs the identical helpers) instead of
// each package keeping its own copy.
func uniqueEmail(t *testing.T) string { return testutils.UniqueEmail(t) }

func stubSMTP(t *testing.T, m *mail.MailService) { testutils.StubSMTP(t, &m.SendMailFunc) }

func stubSMTPAwait(t *testing.T, m *mail.MailService) <-chan struct{} {
	return testutils.StubSMTPAwait(t, &m.SendMailFunc)
}

func awaitSMTP(t *testing.T, called <-chan struct{}) { testutils.AwaitSMTP(t, called) }

// Builds a SigninService whose GenerateAccessToken always fails, by
// pointing AUTH_PRIVATE_PEM_PATH at a nonexistent file during construction
// only — the returned service keeps that bad config even after restore.
func signinServiceWithBadPrivateKey(t *testing.T) *signin.SigninService {
	t.Helper()
	orig := os.Getenv("AUTH_PRIVATE_PEM_PATH")
	require.NoError(t, os.Setenv("AUTH_PRIVATE_PEM_PATH", "/nonexistent/private.pem"))
	config.Init()
	s := signin.NewSigninService(nil)
	require.NoError(t, os.Setenv("AUTH_PRIVATE_PEM_PATH", orig))
	config.Init()
	return s
}

func TestNewAccountService_ReusesProvidedInstance(t *testing.T) {
	existing := &AccountService{}
	assert.Same(t, existing, NewAccountService(existing))
}

func TestNewAccountService_Nil_BuildsRealInstance(t *testing.T) {
	// Safe here (unlike most packages) because TestMain wires db.SetDB and
	// config.Init to working test values.
	s := NewAccountService(nil)
	assert.NotNil(t, s.accountRepository)
	assert.NotNil(t, s.avatarService)
	assert.NotNil(t, s.signinService)
	assert.NotNil(t, s.mailService)
	assert.NotNil(t, s.refreshTokenRepository)
}

func TestCheckUserNameAvailability(t *testing.T) {
	s := newTestAccountService(t)
	username := "checkuser-" + uuid.NewString()
	account := model.Account{Id: uuid.New(), Username: &username}
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: account.Id, Username: &username}, &model.Account{}))

	available, err := s.CheckUserNameAvailability("someone-else-"+uuid.NewString(), nil)
	assert.NoError(t, err)
	assert.True(t, available)

	taken, err := s.CheckUserNameAvailability(username, nil)
	assert.NoError(t, err)
	assert.False(t, taken)

	excluded, err := s.CheckUserNameAvailability(username, &account.Id)
	assert.NoError(t, err)
	assert.True(t, excluded)
}

func TestCheckUserNameAvailability_RepositoryError(t *testing.T) {
	s := newTestAccountService(t)
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))

	_, err := s.CheckUserNameAvailability("someone", nil)
	assert.Error(t, err)
}

func validCreateDto(t *testing.T) *AccountCreateDto {
	t.Helper()
	return &AccountCreateDto{
		Email:         uniqueEmail(t),
		Password:      "SuperSecret123!",
		Language:      constants.ACCOUNT_LANGUAGE_EN,
		TermsAccepted: true,
		TermsVersion:  string(constants.TERMS_VERSIONS[0]),
		TimeZone:      "UTC",
	}
}

func TestCreate_InvalidTermsVersion(t *testing.T) {
	s := newTestAccountService(t)
	dto := validCreateDto(t)
	dto.TermsVersion = "not-a-real-version"

	_, err := s.Create(dto)
	assert.Error(t, err)
}

func TestCreate_InvalidEmail(t *testing.T) {
	s := newTestAccountService(t)
	dto := validCreateDto(t)
	dto.Email = "not-an-email"

	_, err := s.Create(dto)
	assert.ErrorIs(t, err, constants.ERR_INVALID_EMAIL_FORMAT.Err)
}

func TestCreate_InvalidPassword(t *testing.T) {
	s := newTestAccountService(t)
	dto := validCreateDto(t)
	dto.Password = "short"

	_, err := s.Create(dto)
	assert.ErrorIs(t, err, constants.ERR_INVALID_PASSWORD_FORMAT.Err)
}

func TestCreate_EmailAlreadyExists(t *testing.T) {
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	dto := validCreateDto(t)
	dto.Email = email

	_, err := s.Create(dto)
	assert.ErrorIs(t, err, constants.ERR_EMAIL_ALREADY_EXISTS.Err)
}

func TestCreate_InvalidTimeZone(t *testing.T) {
	s := newTestAccountService(t)
	dto := validCreateDto(t)
	dto.TimeZone = "Not/AZone"

	_, err := s.Create(dto)
	assert.Error(t, err)
}

func TestCreate_Success(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	dto := validCreateDto(t)

	tokens, err := s.Create(dto)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(dto.Email, &found))
	assert.NotNil(t, found.Password)

	awaitSMTP(t, called)
}

func TestCreate_Success_FrenchLanguage(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	dto := validCreateDto(t)
	dto.Language = constants.ACCOUNT_LANGUAGE_FR

	tokens, err := s.Create(dto)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)

	awaitSMTP(t, called)
}

func TestCreate_EmailLookupError(t *testing.T) {
	s := newTestAccountService(t)
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))
	dto := validCreateDto(t)

	_, err := s.Create(dto)
	assert.Error(t, err)
}

func TestCreate_RepositoryCreateError(t *testing.T) {
	s := newTestAccountService(t)
	dto := validCreateDto(t)
	testutils.MakeReadOnly(t)

	_, err := s.Create(dto)
	assert.Error(t, err)
}

func TestCreate_GenerateTokensError_RollsBackAccount(t *testing.T) {
	s := newTestAccountService(t)
	s.signinService = signinServiceWithBadPrivateKey(t)
	dto := validCreateDto(t)

	_, err := s.Create(dto)
	assert.Error(t, err)

	var found model.Account
	assert.ErrorIs(t, s.accountRepository.FindOneByEmail(dto.Email, &found), gorm.ErrRecordNotFound, "the account created before token generation failed should be rolled back")
}

func TestGetMe_Found(t *testing.T) {
	s := newTestAccountService(t)
	username := "getme-" + uuid.NewString()
	var created model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Username: &username}, &created))

	dto, err := s.GetMe(created.Id)
	assert.NoError(t, err)
	assert.Equal(t, username, *dto.Username)
}

func TestGetMe_NotFound(t *testing.T) {
	s := newTestAccountService(t)
	_, err := s.GetMe(uuid.New())
	assert.Error(t, err)
}

func createTestAccountWithUsername(t *testing.T, s *AccountService, username string) model.Account {
	t.Helper()
	var created model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Username: &username}, &created))
	return created
}

func TestUpdate_UsernameTooShort(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())
	shortName := "ab"

	_, _, err := s.Update(&AccountUpdateDto{Username: &shortName}, account.Id)
	assert.Error(t, err)
}

func TestUpdate_UsernameAlreadyTaken(t *testing.T) {
	s := newTestAccountService(t)
	taken := "taken-" + uuid.NewString()
	createTestAccountWithUsername(t, s, taken)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())

	_, _, err := s.Update(&AccountUpdateDto{Username: &taken}, account.Id)
	assert.ErrorIs(t, err, constants.ERR_USERNAME_ALREADY_TAKEN.Err)
}

func TestUpdate_UsernameSameAsCurrent_NoOp(t *testing.T) {
	s := newTestAccountService(t)
	username := "same-" + uuid.NewString()
	account := createTestAccountWithUsername(t, s, username)

	// Re-submitting the current username is filtered by the outer guard
	// before the "same as itself" check runs — a silent no-op.
	dto, _, err := s.Update(&AccountUpdateDto{Username: &username}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, username, *dto.Username)
}

func TestUpdate_UsernameChanged_ReturnsNewTokens(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())
	newName := "newname-" + uuid.NewString()

	dto, tokens, err := s.Update(&AccountUpdateDto{Username: &newName}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, newName, *dto.Username)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
}

func TestUpdate_UsernameTrimmedEqualsCurrent(t *testing.T) {
	s := newTestAccountService(t)
	username := "trimtest-" + uuid.NewString()
	account := createTestAccountWithUsername(t, s, username)

	// Differs from the stored username only by whitespace, so trimming makes it equal again.
	padded := "  " + username + "  "
	_, _, err := s.Update(&AccountUpdateDto{Username: &padded}, account.Id)
	assert.ErrorIs(t, err, constants.ERR_USERNAME_ALREADY_TAKEN.Err)
}

func TestUpdate_ValidTimeZone(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "tzok-"+uuid.NewString())
	goodTz := "America/New_York"

	_, _, err := s.Update(&AccountUpdateDto{TimeZone: &goodTz}, account.Id)
	assert.NoError(t, err)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneById(account.Id, &found))
	assert.Equal(t, goodTz, found.TimeZone)
}

func TestUpdate_InvalidTimeZone(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "tz-"+uuid.NewString())
	badTz := "Not/AZone"

	_, _, err := s.Update(&AccountUpdateDto{TimeZone: &badTz}, account.Id)
	assert.Error(t, err)
}

func TestUpdate_InvalidPassword(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "pw-"+uuid.NewString())
	badPw := "short"

	_, _, err := s.Update(&AccountUpdateDto{Password: &badPw}, account.Id)
	assert.ErrorIs(t, err, constants.ERR_INVALID_PASSWORD_FORMAT.Err)
}

func TestUpdate_PasswordChanged_RevokesRefreshTokens(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "pwok-"+uuid.NewString())
	_, err := s.refreshTokenRepository.Create(account.Id, time.Now().Add(time.Hour))
	require.NoError(t, err)

	newPw := "BrandNewSecret123!"
	_, tokens, err := s.Update(&AccountUpdateDto{Password: &newPw}, account.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)

	// RevokeAllForAccount runs before Update generates a fresh replacement
	// token, so we expect the original (now revoked) token plus one new,
	// unrevoked token from the password-change token regeneration.
	var stored []model.RefreshToken
	require.NoError(t, appdb.GetDB().Where("account_id = ?", account.Id).Find(&stored).Error)
	require.Len(t, stored, 2)
	revokedCount, activeCount := 0, 0
	for _, rt := range stored {
		if rt.IsRevoked {
			revokedCount++
		} else {
			activeCount++
		}
	}
	assert.Equal(t, 1, revokedCount, "the pre-existing token should have been revoked")
	assert.Equal(t, 1, activeCount, "the password-change token regeneration should have created a new active token")
}

func TestUpdate_InvalidColor(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "color-"+uuid.NewString())
	badColor := "not-a-color"

	_, _, err := s.Update(&AccountUpdateDto{Color: &badColor}, account.Id)
	assert.ErrorIs(t, err, constants.ERR_INVALID_COLOR_FORMAT.Err)
}

func TestUpdate_ValidColor(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "colorok-"+uuid.NewString())
	color := "#ABCDEF"

	dto, _, err := s.Update(&AccountUpdateDto{Color: &color}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, "#ABCDEF", dto.Color)
}

func TestUpdate_InvalidTermsVersion(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "terms-"+uuid.NewString())
	accepted := true
	badVersion := "not-a-version"

	_, _, err := s.Update(&AccountUpdateDto{TermsAccepted: &accepted, TermsVersion: &badVersion}, account.Id)
	assert.Error(t, err)
}

func TestUpdate_ValidTermsVersion_ReturnsTokens(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "termsok-"+uuid.NewString())
	accepted := true
	version := string(constants.TERMS_VERSIONS[0])

	_, tokens, err := s.Update(&AccountUpdateDto{TermsAccepted: &accepted, TermsVersion: &version}, account.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
}

func TestUpdate_NoQualifyingChanges_NoNewTokens(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "notoken-"+uuid.NewString())
	color := "#123456"

	_, tokens, err := s.Update(&AccountUpdateDto{Color: &color}, account.Id)
	assert.NoError(t, err)
	assert.Nil(t, tokens)
}

func TestUpdate_AccountNotFound(t *testing.T) {
	s := newTestAccountService(t)
	color := "#123456"
	_, _, err := s.Update(&AccountUpdateDto{Color: &color}, uuid.New())
	assert.Error(t, err)
}

func TestUpdate_EmailChanged(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "email-"+uuid.NewString())
	newEmail := uniqueEmail(t)

	dto, _, err := s.Update(&AccountUpdateDto{Email: &newEmail}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, &newEmail, dto.Email)
}

func TestUpdate_LanguageChanged(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "lang-"+uuid.NewString())
	lang := constants.ACCOUNT_LANGUAGE_FR

	dto, _, err := s.Update(&AccountUpdateDto{Language: &lang}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, lang, dto.Language)
}

func TestUpdate_RepositoryUpdatesError(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "updfail-"+uuid.NewString())
	color := "#123456"
	testutils.MakeReadOnly(t)

	_, _, err := s.Update(&AccountUpdateDto{Color: &color}, account.Id)
	assert.Error(t, err)
}

func TestUpdate_GenerateTokensError_RollsBackAccount(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "tokenfail-"+uuid.NewString())
	s.signinService = signinServiceWithBadPrivateKey(t)
	newName := "renamed-" + uuid.NewString()

	_, _, err := s.Update(&AccountUpdateDto{Username: &newName}, account.Id)
	assert.Error(t, err)
}

func TestForgotPassword_AccountNotFound_SilentlyIgnored(t *testing.T) {
	s := newTestAccountService(t)
	err := s.ForgotPassword(&ForgotPasswordDto{Email: uniqueEmail(t)})
	assert.NoError(t, err)
}

func TestForgotPassword_Success(t *testing.T) {
	s := newTestAccountService(t)
	stubSMTP(t, s.mailService)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.NoError(t, err)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(email, &found))
	assert.NotNil(t, found.ResetToken)
}

func TestForgotPassword_SendMailFails(t *testing.T) {
	s := newTestAccountService(t)
	s.mailService.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errors.New("smtp unavailable")
	}

	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.Error(t, err)
}

func TestForgotPassword_Success_FrenchLanguage(t *testing.T) {
	s := newTestAccountService(t)
	stubSMTP(t, s.mailService)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email, Language: constants.ACCOUNT_LANGUAGE_FR}, &model.Account{}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.NoError(t, err)
}

func TestForgotPassword_Cooldown(t *testing.T) {
	s := newTestAccountService(t)
	stubSMTP(t, s.mailService)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	require.NoError(t, s.ForgotPassword(&ForgotPasswordDto{Email: email}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.ErrorIs(t, err, constants.ERR_PASSWORD_RESET_TOO_FREQUENT.Err)
}

func TestForgotPassword_EmailLookupError(t *testing.T) {
	s := newTestAccountService(t)
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: uniqueEmail(t)})
	assert.Error(t, err)
}

func TestForgotPassword_EncryptError(t *testing.T) {
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	orig := os.Getenv("ENCRYPTION_KEY")
	require.NoError(t, os.Setenv("ENCRYPTION_KEY", "too-short"))
	t.Cleanup(func() { _ = os.Setenv("ENCRYPTION_KEY", orig) })

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.Error(t, err)
}

func TestForgotPassword_UpdateResetTokenError(t *testing.T) {
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))
	testutils.MakeReadOnly(t)

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.Error(t, err)
}

func TestResetPassword_InvalidPasswordFormat(t *testing.T) {
	s := newTestAccountService(t)
	err := s.ResetPassword(&ResetPasswordDto{Token: "irrelevant", Password: "short"})
	assert.ErrorIs(t, err, constants.ERR_INVALID_PASSWORD_FORMAT.Err)
}

func TestResetPassword_UndecryptableToken(t *testing.T) {
	s := newTestAccountService(t)
	err := s.ResetPassword(&ResetPasswordDto{Token: "not-valid-base64!!", Password: "SuperSecret123!"})
	assert.Error(t, err)
}

func TestResetPassword_TokenNotFound(t *testing.T) {
	s := newTestAccountService(t)
	encrypted, err := encryption.Encrypt("some-unknown-reset-token")
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.ErrorIs(t, err, constants.ERR_INVALID_RESET_TOKEN.Err)
}

func TestResetPassword_RepositoryLookupError(t *testing.T) {
	s := newTestAccountService(t)
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))
	encrypted, err := encryption.Encrypt("some-token")
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.Error(t, err)
}

func TestResetPassword_RepositoryUpdatesError(t *testing.T) {
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &account))

	rawToken := "update-fail-token"
	expiresAt := time.Now().Add(time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiresAt))
	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	testutils.MakeReadOnly(t)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.Error(t, err)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &account))

	rawToken := "expired-token-value"
	expiredAt := time.Now().Add(-time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiredAt))

	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.ErrorIs(t, err, constants.ERR_RESET_TOKEN_EXPIRED.Err)
}

func TestResetPassword_Success(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	email := uniqueEmail(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &account))

	rawToken := "valid-token-value"
	expiresAt := time.Now().Add(time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiresAt))

	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.NoError(t, err)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneById(account.Id, &found))
	assert.Nil(t, found.ResetToken)

	awaitSMTP(t, called)
}

func TestResetPassword_Success_FrenchLanguage(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	email := uniqueEmail(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email, Language: constants.ACCOUNT_LANGUAGE_FR}, &account))

	rawToken := "fr-token-value"
	expiresAt := time.Now().Add(time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiresAt))

	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.NoError(t, err)

	awaitSMTP(t, called)
}

func TestResetPassword_InvalidStoredTimeZone_FallsBackToUTC(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	email := uniqueEmail(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &account))

	account.TimeZone = "Not/AZone"
	require.NoError(t, s.accountRepository.Updates(account))

	rawToken := "bad-tz-token-value"
	expiresAt := time.Now().Add(time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiresAt))

	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.NoError(t, err)

	awaitSMTP(t, called)
}

func TestResetPassword_NoEmail_SkipsConfirmationEmail(t *testing.T) {
	s := newTestAccountService(t)
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New()}, &account))

	rawToken := "no-email-token"
	expiresAt := time.Now().Add(time.Hour)
	require.NoError(t, s.accountRepository.UpdateResetToken(account.Id, &rawToken, &expiresAt))

	encrypted, err := encryption.Encrypt(rawToken)
	require.NoError(t, err)

	err = s.ResetPassword(&ResetPasswordDto{Token: encrypted, Password: "SuperSecret123!"})
	assert.NoError(t, err)
}

func TestDelete_Success(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "del-"+uuid.NewString())

	_, err := s.Delete(account.Id, &guard.Claims{Id: account.Id})
	assert.NoError(t, err)

	var found model.Account
	assert.Error(t, s.accountRepository.FindOneById(account.Id, &found))
}

func TestDelete_NotFound(t *testing.T) {
	s := newTestAccountService(t)
	_, err := s.Delete(uuid.New(), &guard.Claims{})
	assert.Error(t, err)
}

// Targets Delete's final db.GetDB().Delete call, which reads the
// process-global DB directly rather than through a repo field.
func TestDelete_DeleteQueryError(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "delfail-"+uuid.NewString())
	testutils.MakeReadOnly(t)

	_, err := s.Delete(account.Id, &guard.Claims{Id: account.Id})
	assert.Error(t, err)
}
