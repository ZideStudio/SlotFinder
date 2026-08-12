package provider

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/config"
	appdb "app/db"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/account"
	"app/pkg/mail"
	"app/pkg/signin"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

		discordTokenURL:    constants.PROVIDER_DISCORD_TOKEN_URL,
		discordUserInfoURL: constants.PROVIDER_DISCORD_USERINFO_URL,

		googleTokenURL:    constants.PROVIDER_GOOGLE_TOKEN_URL,
		googleUserInfoURL: constants.PROVIDER_GOOGLE_USERINFO_URL,

		githubTokenURL:     constants.PROVIDER_GITHUB_TOKEN_URL,
		githubUserInfoURL:  constants.PROVIDER_GITHUB_USERINFO_URL,
		githubUserEmailURL: constants.PROVIDER_GITHUB_USEREMAIL_URL,
	}
}

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s@example.com", uuid.NewString())
}

func stubSMTPAwait(t *testing.T, m *mail.MailService) <-chan struct{} {
	t.Helper()
	called := make(chan struct{}, 1)
	m.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called <- struct{}{}
		return nil
	}
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

// closedRepoDB returns a gorm.DB whose underlying connection is already
// closed, so any query through it fails immediately — used to force a
// repository-level error independently of the other repositories on the
// service, which keep using the real shared test DB.
func closedRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := database.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
	return database
}

// setQueryOnly flips the shared test DB read-only for the duration of the
// test, so prior SELECTs in the code path under test keep working while the
// next write fails. Restored via t.Cleanup.
func setQueryOnly(t *testing.T, on bool) {
	t.Helper()
	sqlDB, err := appdb.GetDB().DB()
	require.NoError(t, err)
	value := "OFF"
	if on {
		value = "ON"
	}
	_, err = sqlDB.Exec("PRAGMA query_only = " + value)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = sqlDB.Exec("PRAGMA query_only = OFF") })
}

// signinServiceWithBadPrivateKey builds a signin.SigninService whose
// GenerateAccessToken always fails, by pointing AUTH_PRIVATE_PEM_PATH at a
// nonexistent file for the duration of construction only.
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

func TestNewProviderService_ReusesProvidedInstance(t *testing.T) {
	existing := &ProviderService{}
	assert.Same(t, existing, NewProviderService(existing))
}

func TestNewProviderService_Nil_BuildsRealDependencies(t *testing.T) {
	s := NewProviderService(nil)
	assert.NotNil(t, s.accountProvidersRepository)
	assert.NotNil(t, s.accountRepository)
	assert.NotNil(t, s.signinService)
	assert.NotNil(t, s.accountService)
	assert.NotNil(t, s.avatarService)
	assert.NotNil(t, s.mailService)
	assert.NotEmpty(t, s.discordTokenURL)
	assert.NotEmpty(t, s.googleTokenURL)
	assert.NotEmpty(t, s.githubTokenURL)
}

// tokenOnlyServer serves a valid "/token" response so the second-hop URL can be tested in isolation.
func tokenOnlyServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	}))
	t.Cleanup(server.Close)
	return server
}

func TestGetGoogleUserInfo_UserInfoNetworkError(t *testing.T) {
	tokenServer := tokenOnlyServer(t)
	s := newTestProviderService(t)
	s.googleTokenURL = tokenServer.URL
	s.googleUserInfoURL = closedServerURL(t)

	_, err := s.getGoogleUserInfo("code")
	assert.Error(t, err)
}

func TestGetDiscordUserInfo_UserInfoNetworkError(t *testing.T) {
	tokenServer := tokenOnlyServer(t)
	s := newTestProviderService(t)
	s.discordTokenURL = tokenServer.URL
	s.discordUserInfoURL = closedServerURL(t)

	_, err := s.getDiscordUserInfo("code")
	assert.Error(t, err)
}

func TestGetGithubUserInfo_UserInfoNetworkError(t *testing.T) {
	tokenServer := tokenOnlyServer(t)
	s := newTestProviderService(t)
	s.githubTokenURL = tokenServer.URL
	s.githubUserInfoURL = closedServerURL(t)

	_, err := s.getGithubUserInfo("code")
	assert.Error(t, err)
}

func TestGetGithubUserInfo_UserInfoNonSuccessStatus(t *testing.T) {
	tokenServer := tokenOnlyServer(t)
	failureServer := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.githubTokenURL = tokenServer.URL
	s.githubUserInfoURL = failureServer.URL

	_, err := s.getGithubUserInfo("code")
	assert.Error(t, err)
}

func TestGetGithubUserInfo_EmailsNetworkError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": 1, "login": "u", "avatar_url": ""})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL = server.URL+"/token", server.URL+"/user"
	s.githubUserEmailURL = closedServerURL(t)

	_, err := s.getGithubUserInfo("code")
	assert.Error(t, err)
}

func TestGetGithubUserInfo_EmailsNonSuccessStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": 1, "login": "u", "avatar_url": ""})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	failureServer := newOAuthFailureServer(t)

	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL = server.URL+"/token", server.URL+"/user"
	s.githubUserEmailURL = failureServer.URL

	_, err := s.getGithubUserInfo("code")
	assert.Error(t, err)
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

func TestGetProviderUrl_EncryptError(t *testing.T) {
	s := newTestProviderService(t)
	orig := os.Getenv("ENCRYPTION_KEY")
	require.NoError(t, os.Setenv("ENCRYPTION_KEY", "too-short"))
	t.Cleanup(func() { _ = os.Setenv("ENCRYPTION_KEY", orig) })

	_, err := s.GetProviderUrl("google", "/dashboard", nil)
	assert.Error(t, err)
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

func TestCreateProviderAccount_EmailLookupError(t *testing.T) {
	s := newTestProviderService(t)
	s.accountRepository = repository.NewAccountRepository(closedRepoDB(t))
	email := uniqueEmail(t)

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")
	assert.Error(t, err)
}

func TestCreateProviderAccount_ProviderLookupError(t *testing.T) {
	s := newTestProviderService(t)
	s.accountProvidersRepository = repository.NewAccountProvidersRepository(closedRepoDB(t))
	email := uniqueEmail(t)

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")
	assert.Error(t, err)
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

func TestCreateProviderAccount_ReLogin_GenerateTokensError(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	providerId := "google-" + uuid.NewString()
	account := createAccountWithProvider(t, s, email, constants.PROVIDER_GOOGLE, providerId)
	s.signinService = signinServiceWithBadPrivateKey(t)

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: providerId, Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, account.Id.String())
	assert.Error(t, err)
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

// TestCreateProviderAccount_ProviderLinkedToAnotherAccount_NoEmailConflict
// uses a fresh (unregistered) OAuth email, so the email-conflict check
// (ERR_EMAIL_ALREADY_EXISTS) never fires and the "ALREADY_EXISTS" error
// genuinely comes from the provider-ownership check below it — unlike
// TestCreateProviderAccount_ProviderLinkedToAnotherAccount above, whose
// OAuth email collides with an existing account and is rejected earlier.
func TestCreateProviderAccount_ProviderLinkedToAnotherAccount_NoEmailConflict(t *testing.T) {
	s := newTestProviderService(t)
	ownerEmail := uniqueEmail(t)
	providerId := "google-" + uuid.NewString()
	createAccountWithProvider(t, s, ownerEmail, constants.PROVIDER_GOOGLE, providerId)

	var other model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New()}, &other))

	freshEmail := uniqueEmail(t)
	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: providerId, Email: &freshEmail},
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
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = googleServer.URL+"/token", googleServer.URL+"/userinfo"

	_, err := s.ProviderCallback("google", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Google_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = url+"/token", url+"/userinfo"

	_, err := s.ProviderCallback("google", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Google_UserInfoNonSuccessStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Google_NewAccount(t *testing.T) {
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)

	// The mock userinfo response has an empty "picture" field, so
	// avatarService.UploadAvatar(&"", nil, ...) fails fast (invalid URL) and
	// ProviderCallback falls back to GetGravatarURL — no imgbb stub needed.
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	called := stubSMTPAwait(t, s.mailService)

	tokens, err := s.ProviderCallback("google", "good-code", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(email, &found))

	awaitSMTP(t, called)
}

func TestProviderCallback_CreateProviderAccountError(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	createAccountWithProvider(t, s, email, constants.PROVIDER_GITHUB, "gh-"+uuid.NewString())

	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "someuser", email)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.ErrorIs(t, err, constants.ERR_EMAIL_ALREADY_EXISTS.Err)
}

func TestProviderCallback_NewAccount_RepositoryCreateError(t *testing.T) {
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	setQueryOnly(t, true)

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_NewAccount_GenerateTokensError(t *testing.T) {
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)
	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	s.signinService = signinServiceWithBadPrivateKey(t)

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.Error(t, err)

	var found model.Account
	assert.ErrorIs(t, s.accountRepository.FindOneByEmail(email, &found), gorm.ErrRecordNotFound, "the account created before token generation failed should be rolled back")
}

func TestProviderCallback_Discord_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = url+"/token", url+"/userinfo"

	_, err := s.ProviderCallback("discord", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_TokenNonSuccessStatus(t *testing.T) {
	server := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("discord", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_UserInfoNonSuccessStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_NewAccount_MissingUsername(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-id-"+uuid.NewString(), "", uniqueEmail(t))
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL, s.githubUserEmailURL = url+"/token", url+"/user", url+"/emails"

	_, err := s.ProviderCallback("github", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_TokenNonSuccessStatus(t *testing.T) {
	server := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL, s.githubUserEmailURL = server.URL+"/token", server.URL+"/user", server.URL+"/emails"

	_, err := s.ProviderCallback("github", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_NoPrimaryEmail_FallsBackToFirst(t *testing.T) {
	firstEmail := uniqueEmail(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": 999, "login": "no-primary-user", "avatar_url": ""})
	})
	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []map[string]any{
			{"email": firstEmail, "primary": false},
			{"email": uniqueEmail(t), "primary": false},
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL, s.githubUserEmailURL = server.URL+"/token", server.URL+"/user", server.URL+"/emails"

	resp, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: func() ProviderAccount {
			userInfo, err := s.getGithubUserInfo("good-code")
			require.NoError(t, err)
			return userInfo
		}(),
		Provider: constants.PROVIDER_GITHUB,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, resp.Account)
	assert.Equal(t, firstEmail, *resp.Account.Email)
}

func TestProviderCallback_Github_NoEmails_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": 998, "login": "no-email-user", "avatar_url": ""})
	})
	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []map[string]any{})
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	s := newTestProviderService(t)
	s.githubTokenURL, s.githubUserInfoURL, s.githubUserEmailURL = server.URL+"/token", server.URL+"/user", server.URL+"/emails"

	_, err := s.ProviderCallback("github", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_MalformedUserId(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("discord", "good-code", "not-a-uuid")
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_UnknownUserId(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"

	_, err := s.ProviderCallback("discord", "good-code", uuid.New().String())
	assert.Error(t, err)
}

// clearNilAccountIdProviders removes any AccountProvider rows left over with
// a zero AccountId by other "brand-new-link" tests, which would otherwise
// collide with this test's own on the (account_id, provider) unique index.
func clearNilAccountIdProviders(t *testing.T) {
	t.Helper()
	require.NoError(t, appdb.GetDB().Where("account_id = ?", uuid.Nil).Delete(&model.AccountProvider{}).Error)
}

func TestProviderCallback_LinkProvider_AccountLookupError(t *testing.T) {
	clearNilAccountIdProviders(t)
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	s.accountRepository = repository.NewAccountRepository(closedRepoDB(t))

	_, err := s.ProviderCallback("discord", "good-code", uuid.New().String())
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_GenerateTokensError(t *testing.T) {
	clearNilAccountIdProviders(t)
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.discordTokenURL, s.discordUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	loggedInUser := createTestAccountForLink(t, s)
	s.signinService = signinServiceWithBadPrivateKey(t)

	_, err := s.ProviderCallback("discord", "good-code", loggedInUser.String())
	assert.Error(t, err)
}

func createTestAccountForLink(t *testing.T, s *ProviderService) uuid.UUID {
	t.Helper()
	var account model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New()}, &account))
	return account.Id
}

func buildTestAvatarServer(t *testing.T) *httptest.Server {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	body := buf.Bytes()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestProviderCallback_Google_NewAccount_FetchesAvatarFromProviderURL(t *testing.T) {
	email := uniqueEmail(t)
	avatarServer := buildTestAvatarServer(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"sub": "google-avatar-" + uuid.NewString(), "name": "avataruser", "email": email, "picture": avatarServer.URL})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	s := newTestProviderService(t)
	s.googleTokenURL, s.googleUserInfoURL = server.URL+"/token", server.URL+"/userinfo"
	called := stubSMTPAwait(t, s.mailService)

	tokens, err := s.ProviderCallback("google", "good-code", "")
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)

	var found model.Account
	require.NoError(t, s.accountRepository.FindOneByEmail(email, &found))
	assert.NotEmpty(t, found.AvatarData)
	assert.Contains(t, found.AvatarUrl, "/api/v1/account/")

	awaitSMTP(t, called)
}

func TestProviderCallback_Github_ExistingAccountLogin(t *testing.T) {
	email := uniqueEmail(t)
	githubId := "12345"
	s := newTestProviderService(t)
	createAccountWithProvider(t, s, email, constants.PROVIDER_GITHUB, githubId)

	server := newGithubOAuthServer(t, 12345, "githubuser", email)
	s.githubTokenURL, s.githubUserInfoURL, s.githubUserEmailURL = server.URL+"/token", server.URL+"/user", server.URL+"/emails"

	tokens, err := s.ProviderCallback("github", "good-code", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
}
