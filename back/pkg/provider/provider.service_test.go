package provider

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/config"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/account"
	"app/pkg/mail"
	"app/pkg/signin"
	"fmt"
	"net/smtp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: Do not call NewProviderService(nil) in most tests — build the struct
// directly instead (same rationale documented in pkg/account tests). Nested
// NewXService(nil) calls are safe here because TestMain wires db.SetDB and
// config.Init.
func newTestProviderService(t *testing.T) *ProviderService {
	t.Helper()
	return &ProviderService{
		accountProvidersRepository: repository.NewAccountProvidersRepository(nil),
		accountRepository:          repository.NewAccountRepository(nil),
		signinService:              signin.NewSigninService(nil),
		accountService:             account.NewAccountService(nil),
		avatarService:              account.NewAvatarService(nil),
		mailService:                mail.NewMailService(nil),
		config:                     config.GetConfig(),
	}
}

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s@example.com", uuid.NewString())
}

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

func TestNewProviderService_ReusesProvidedInstance(t *testing.T) {
	existing := &ProviderService{}
	assert.Same(t, existing, NewProviderService(existing))
}

func TestParseProvider(t *testing.T) {
	s := &ProviderService{}

	google, err := s.parseProvider("google")
	assert.NoError(t, err)
	assert.Equal(t, constants.PROVIDER_GOOGLE, google)

	discord, err := s.parseProvider("discord")
	assert.NoError(t, err)
	assert.Equal(t, constants.PROVIDER_DISCORD, discord)

	github, err := s.parseProvider("github")
	assert.NoError(t, err)
	assert.Equal(t, constants.PROVIDER_GITHUB, github)

	_, err = s.parseProvider("unknown")
	assert.Error(t, err)
}

func TestGetProviderUrl_EachProvider(t *testing.T) {
	s := newTestProviderService(t)

	googleUrl, err := s.GetProviderUrl("google", "/dashboard", nil)
	assert.NoError(t, err)
	assert.Contains(t, googleUrl, "https://accounts.google.com")

	discordUrl, err := s.GetProviderUrl("discord", "/dashboard", nil)
	assert.NoError(t, err)
	assert.Contains(t, discordUrl, "https://discord.com/oauth2/authorize")

	githubUrl, err := s.GetProviderUrl("github", "/dashboard", nil)
	assert.NoError(t, err)
	assert.Contains(t, githubUrl, "https://github.com/login/oauth/authorize")
}

func TestGetProviderUrl_InvalidProvider(t *testing.T) {
	s := newTestProviderService(t)
	_, err := s.GetProviderUrl("not-a-provider", "", nil)
	assert.Error(t, err)
}

func TestGetProviderUrl_WithAuthenticatedUser(t *testing.T) {
	s := newTestProviderService(t)
	url, err := s.GetProviderUrl("google", "/dashboard", &guard.Claims{Id: uuid.New()})
	assert.NoError(t, err)
	assert.NotEmpty(t, url)
}

func createAccountWithProvider(t *testing.T, s *ProviderService, email string, provider constants.Provider, providerId string) model.Account {
	t.Helper()
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{
		Id:    uuid.New(),
		Email: &email,
		Providers: []model.AccountProvider{
			{Provider: provider, Id: providerId},
		},
	}, &account))
	return account
}

func TestCreateProviderAccount_EmailAlreadyExists_DifferentProvider(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	createAccountWithProvider(t, s, email, constants.PROVIDER_GITHUB, "gh-"+uuid.NewString())

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")

	assert.ErrorIs(t, err, constants.ERR_EMAIL_ALREADY_EXISTS.Err)
}

func TestCreateProviderAccount_NewAccount(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	username := "newuser-" + uuid.NewString()

	resp, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &email, Username: username},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")

	assert.NoError(t, err)
	assert.NotNil(t, resp.Account)
	assert.Equal(t, &email, resp.Account.Email)
	assert.Nil(t, resp.Jwt)
	assert.Nil(t, resp.AccountProvider)
}

func TestCreateProviderAccount_ExistingProviderLogin(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	providerId := "google-" + uuid.NewString()
	createAccountWithProvider(t, s, email, constants.PROVIDER_GOOGLE, providerId)

	resp, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: providerId, Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")

	assert.NoError(t, err)
	require.NotNil(t, resp.Jwt)
	assert.NotEmpty(t, resp.Jwt.AccessToken)
}

func TestCreateProviderAccount_LinkNewProviderToLoggedInAccount(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	account := createAccountWithProvider(t, s, email, constants.PROVIDER_GITHUB, "gh-"+uuid.NewString())

	// The user is logged in (authUserId set) and links a Google account
	// whose OAuth email isn't tied to any existing account (the email-based
	// conflict check only looks at whether *that* email is already
	// registered, independent of the caller's own account/email).
	newProviderEmail := uniqueEmail(t)
	resp, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &newProviderEmail},
		Provider:        constants.PROVIDER_GOOGLE,
	}, account.Id.String())

	assert.NoError(t, err)
	require.NotNil(t, resp.AccountProvider)
	assert.Equal(t, constants.PROVIDER_GOOGLE, resp.AccountProvider.Provider)
}

func TestCreateProviderAccount_ReLoginViaOwnLinkedProvider(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	providerId := "google-" + uuid.NewString()
	account := createAccountWithProvider(t, s, email, constants.PROVIDER_GOOGLE, providerId)

	resp, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: providerId, Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, account.Id.String())

	assert.NoError(t, err)
	require.NotNil(t, resp.Jwt)
}

func TestCreateProviderAccount_ProviderLinkedToAnotherAccount(t *testing.T) {
	s := newTestProviderService(t)
	ownerEmail := uniqueEmail(t)
	providerId := "google-" + uuid.NewString()
	createAccountWithProvider(t, s, ownerEmail, constants.PROVIDER_GOOGLE, providerId)

	otherEmail := uniqueEmail(t)
	var other model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &otherEmail}, &other))

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: providerId, Email: &otherEmail},
		Provider:        constants.PROVIDER_GOOGLE,
	}, other.Id.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ALREADY_EXISTS")
}

func TestProviderCallback_InvalidProvider(t *testing.T) {
	s := newTestProviderService(t)
	_, err := s.ProviderCallback("not-a-provider", "code", "")
	assert.Error(t, err)
}

func TestProviderCallback_GoogleOAuthFailure(t *testing.T) {
	googleServer := newOAuthFailureServer(t)
	restoreGoogle := setGoogleURLs(googleServer.URL+"/token", googleServer.URL+"/userinfo")
	defer restoreGoogle()

	s := newTestProviderService(t)
	_, err := s.ProviderCallback("google", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Google_NewAccount(t *testing.T) {
	called := stubSMTPAwait(t)
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)
	restore := setGoogleURLs(server.URL+"/token", server.URL+"/userinfo")
	defer restore()

	// The mock userinfo response has an empty "picture" field, so
	// avatarService.UploadAvatar(&"", nil, ...) fails fast (invalid URL) and
	// ProviderCallback falls back to GetGravatarURL — no imgbb stub needed.
	s := newTestProviderService(t)
	tokens, err := s.ProviderCallback("google", "good-code", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(email, &found))

	awaitSMTP(t, called)
}

func TestProviderCallback_Discord_NewAccount_MissingUsername(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-id-"+uuid.NewString(), "", uniqueEmail(t))
	restore := setDiscordURLs(server.URL+"/token", server.URL+"/userinfo")
	defer restore()

	s := newTestProviderService(t)
	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_ExistingAccountLogin(t *testing.T) {
	email := uniqueEmail(t)
	githubId := "12345"
	s := newTestProviderService(t)
	createAccountWithProvider(t, s, email, constants.PROVIDER_GITHUB, githubId)

	server := newGithubOAuthServer(t, 12345, "githubuser", email)
	restore := setGithubURLs(server.URL+"/token", server.URL+"/user", server.URL+"/emails")
	defer restore()

	tokens, err := s.ProviderCallback("github", "good-code", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
}
