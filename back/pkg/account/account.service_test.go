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
	"fmt"
	"net/smtp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: Do not call NewAccountService(nil) in tests — build the struct
// directly. Nested services (signinService/avatarService) are still built
// via their own NewXService(nil), which is safe here because TestMain wires
// db.SetDB to an in-memory sqlite DB and config.Init to real (test) values;
// see the note on TestMain in main_test.go.
func newTestAccountService(t *testing.T) *AccountService {
	t.Helper()
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

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s@example.com", uuid.NewString())
}

// stubSMTP stubs mail.SmtpSendFunc for the duration of the test. Only safe
// for code paths that call SendMail synchronously (e.g. ForgotPassword) —
// for paths that send via `go s.mailService.SendMail(...)` (Create,
// ResetPassword), use stubSMTPAwait instead so the test doesn't restore the
// stub while that goroutine might still be reading it (data race).
func stubSMTP(t *testing.T) {
	t.Helper()
	original := mail.SmtpSendFunc
	mail.SmtpSendFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error { return nil }
	t.Cleanup(func() { mail.SmtpSendFunc = original })
}

// stubSMTPAwait stubs mail.SmtpSendFunc and returns a channel that receives
// once the stub has actually been invoked, so callers can wait for an
// asynchronously-spawned SendMail goroutine to run before the test (and its
// cleanup restoring the original SmtpSendFunc) returns.
func stubSMTPAwait(t *testing.T) <-chan struct{} {
	t.Helper()
	called := make(chan struct{}, 1)
	original := mail.SmtpSendFunc
	mail.SmtpSendFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called <- struct{}{}
		return nil
	}
	t.Cleanup(func() { mail.SmtpSendFunc = original })
	return called
}

func awaitSMTP(t *testing.T, called <-chan struct{}) {
	t.Helper()
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("expected the async SendMail goroutine to run")
	}
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
	account := model.Account{Id: uuid.New(), UserName: &username}
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: account.Id, UserName: &username}, &model.Account{}))

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
	called := stubSMTPAwait(t)
	s := newTestAccountService(t)
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

func TestGetMe_Found(t *testing.T) {
	s := newTestAccountService(t)
	username := "getme-" + uuid.NewString()
	var created model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), UserName: &username}, &created))

	dto, err := s.GetMe(created.Id)
	assert.NoError(t, err)
	assert.Equal(t, username, *dto.UserName)
}

func TestGetMe_NotFound(t *testing.T) {
	s := newTestAccountService(t)
	_, err := s.GetMe(uuid.New())
	assert.Error(t, err)
}

func createTestAccountWithUsername(t *testing.T, s *AccountService, username string) model.Account {
	t.Helper()
	var created model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), UserName: &username}, &created))
	return created
}

func TestUpdate_UsernameTooShort(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())
	shortName := "ab"

	_, _, err := s.Update(&AccountUpdateDto{UserName: &shortName}, account.Id)
	assert.Error(t, err)
}

func TestUpdate_UsernameAlreadyTaken(t *testing.T) {
	s := newTestAccountService(t)
	taken := "taken-" + uuid.NewString()
	createTestAccountWithUsername(t, s, taken)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())

	_, _, err := s.Update(&AccountUpdateDto{UserName: &taken}, account.Id)
	assert.ErrorIs(t, err, constants.ERR_USERNAME_ALREADY_TAKEN.Err)
}

func TestUpdate_UsernameSameAsCurrent_NoOp(t *testing.T) {
	s := newTestAccountService(t)
	username := "same-" + uuid.NewString()
	account := createTestAccountWithUsername(t, s, username)

	// Re-submitting the current username is filtered out by the outer
	// `*dto.UserName != *account.UserName` guard before the "same as
	// itself" check further down ever runs — so this is a silent no-op,
	// not ERR_USERNAME_ALREADY_TAKEN.
	dto, _, err := s.Update(&AccountUpdateDto{UserName: &username}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, username, *dto.UserName)
}

func TestUpdate_UsernameChanged_ReturnsNewTokens(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "orig-"+uuid.NewString())
	newName := "newname-" + uuid.NewString()

	dto, tokens, err := s.Update(&AccountUpdateDto{UserName: &newName}, account.Id)
	assert.NoError(t, err)
	assert.Equal(t, newName, *dto.UserName)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
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

func TestForgotPassword_AccountNotFound_SilentlyIgnored(t *testing.T) {
	s := newTestAccountService(t)
	err := s.ForgotPassword(&ForgotPasswordDto{Email: uniqueEmail(t)})
	assert.NoError(t, err)
}

func TestForgotPassword_Success(t *testing.T) {
	stubSMTP(t)
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.NoError(t, err)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(email, &found))
	assert.NotNil(t, found.ResetToken)
}

func TestForgotPassword_Cooldown(t *testing.T) {
	stubSMTP(t)
	s := newTestAccountService(t)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))

	require.NoError(t, s.ForgotPassword(&ForgotPasswordDto{Email: email}))

	err := s.ForgotPassword(&ForgotPasswordDto{Email: email})
	assert.ErrorIs(t, err, constants.ERR_PASSWORD_RESET_TOO_FREQUENT.Err)
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
	called := stubSMTPAwait(t)
	s := newTestAccountService(t)
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
