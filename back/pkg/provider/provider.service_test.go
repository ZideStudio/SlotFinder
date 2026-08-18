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
	"app/testutils"
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// NOTE: Do not call NewProviderService(nil) in most tests — build the struct
// directly instead (same rationale as pkg/account's tests).
func newTestProviderService(t *testing.T) *ProviderService {
	t.Helper()
	testutils.TestDB(t)
	return &ProviderService{
		accountProvidersRepository: repository.NewAccountProvidersRepository(nil),
		accountRepository:          repository.NewAccountRepository(nil),
		signinService:              signin.NewSigninService(nil),
		accountService:             account.NewAccountService(nil),
		avatarService:              account.NewAvatarService(nil),
		mailService:                mail.NewMailService(nil),
		config:                     config.GetConfig(),
		httpClient:                 resty.New(),
	}
}

// uniqueEmail, stubSMTPAwait, and awaitSMTP delegate to testutils (shared
// with pkg/account, which needs the identical helpers) instead of each
// package keeping its own copy.
func uniqueEmail(t *testing.T) string { return testutils.UniqueEmail(t) }

func stubSMTPAwait(t *testing.T, m *mail.MailService) <-chan struct{} {
	return testutils.StubSMTPAwait(t, &m.SendMailFunc)
}

func awaitSMTP(t *testing.T, called <-chan struct{}) { testutils.AwaitSMTP(t, called) }

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
	assert.NotNil(t, s.httpClient)
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

func TestGetUserInfo_UserInfoNetworkError(t *testing.T) {
	cases := []struct {
		name        string
		tokenURL    string
		userInfoURL string
		call        func(s *ProviderService, code string) (ProviderAccount, error)
	}{
		{
			name:        "google",
			tokenURL:    realGoogleTokenURL,
			userInfoURL: realGoogleUserInfoURL,
			call:        (*ProviderService).getGoogleUserInfo,
		},
		{
			name:        "discord",
			tokenURL:    realDiscordTokenURL,
			userInfoURL: realDiscordUserInfoURL,
			call:        (*ProviderService).getDiscordUserInfo,
		},
		{
			name:        "github",
			tokenURL:    realGithubTokenURL,
			userInfoURL: realGithubUserInfoURL,
			call:        (*ProviderService).getGithubUserInfo,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokenServer := tokenOnlyServer(t)
			s := newTestProviderService(t)
			s.httpClient = oauthRedirectClient(t, map[string]string{
				tc.tokenURL:    tokenServer.URL,
				tc.userInfoURL: closedServerURL(t),
			})

			_, err := tc.call(s, "code")
			assert.Error(t, err)
		})
	}
}

func TestGetGithubUserInfo_UserInfoNonSuccessStatus(t *testing.T) {
	tokenServer := tokenOnlyServer(t)
	failureServer := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:    tokenServer.URL,
		realGithubUserInfoURL: failureServer.URL,
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: closedServerURL(t),
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: failureServer.URL,
	})

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
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))
	email := uniqueEmail(t)

	_, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: ProviderAccount{Id: "google-" + uuid.NewString(), Email: &email},
		Provider:        constants.PROVIDER_GOOGLE,
	}, "")
	assert.Error(t, err)
}

func TestCreateProviderAccount_ProviderLookupError(t *testing.T) {
	s := newTestProviderService(t)
	s.accountProvidersRepository = repository.NewAccountProvidersRepository(testutils.ClosedDB(t))
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

	// Logged-in user links a Google account whose OAuth email isn't tied
	// to any existing account.
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

// Uses a fresh OAuth email so the email-conflict check never fires, and the
// "ALREADY_EXISTS" error genuinely comes from the provider-ownership check.
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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    googleServer.URL + "/token",
		realGoogleUserInfoURL: googleServer.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("google", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Google_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    url + "/token",
		realGoogleUserInfoURL: url + "/userinfo",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})
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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.ErrorIs(t, err, constants.ERR_EMAIL_ALREADY_EXISTS.Err)
}

func TestProviderCallback_NewAccount_RepositoryCreateError(t *testing.T) {
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})
	testutils.MakeReadOnly(t)

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_NewAccount_GenerateTokensError(t *testing.T) {
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "newgoogleuser", email)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})
	s.signinService = signinServiceWithBadPrivateKey(t)

	_, err := s.ProviderCallback("google", "good-code", "")
	assert.Error(t, err)

	var found model.Account
	assert.ErrorIs(t, s.accountRepository.FindOneByEmail(email, &found), gorm.ErrRecordNotFound, "the account created before token generation failed should be rolled back")
}

func TestProviderCallback_Discord_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    url + "/token",
		realDiscordUserInfoURL: url + "/userinfo",
	})

	_, err := s.ProviderCallback("discord", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_TokenNonSuccessStatus(t *testing.T) {
	server := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_NewAccount_MissingUsername(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-id-"+uuid.NewString(), "", uniqueEmail(t))
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Discord_NewAccount_MissingEmail(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-id-"+uuid.NewString(), "newdiscorduser", "")
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("discord", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_TokenNetworkError(t *testing.T) {
	url := closedServerURL(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     url + "/token",
		realGithubUserInfoURL:  url + "/user",
		realGithubUserEmailURL: url + "/emails",
	})

	_, err := s.ProviderCallback("github", "bad-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_Github_TokenNonSuccessStatus(t *testing.T) {
	server := newOAuthFailureServer(t)
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: server.URL + "/emails",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: server.URL + "/emails",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: server.URL + "/emails",
	})

	_, err := s.ProviderCallback("github", "good-code", "")
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_MalformedUserId(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

	_, err := s.ProviderCallback("discord", "good-code", "not-a-uuid")
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_UnknownUserId(t *testing.T) {
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})

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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})
	s.accountRepository = repository.NewAccountRepository(testutils.ClosedDB(t))

	_, err := s.ProviderCallback("discord", "good-code", uuid.New().String())
	assert.Error(t, err)
}

func TestProviderCallback_LinkProvider_GenerateTokensError(t *testing.T) {
	clearNilAccountIdProviders(t)
	server := newDiscordOAuthServer(t, "discord-link-"+uuid.NewString(), "linkuser", uniqueEmail(t))
	s := newTestProviderService(t)
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realDiscordTokenURL:    server.URL + "/token",
		realDiscordUserInfoURL: server.URL + "/userinfo",
	})
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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGoogleTokenURL:    server.URL + "/token",
		realGoogleUserInfoURL: server.URL + "/userinfo",
	})
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
	s.httpClient = oauthRedirectClient(t, map[string]string{
		realGithubTokenURL:     server.URL + "/token",
		realGithubUserInfoURL:  server.URL + "/user",
		realGithubUserEmailURL: server.URL + "/emails",
	})

	tokens, err := s.ProviderCallback("github", "good-code", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
}
