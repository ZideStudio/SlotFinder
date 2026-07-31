package provider

import (
	"app/commons/encryption"
	model "app/db/models"
	"app/db/repository"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newProviderTestContext(providerParam, query string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	target := "/?" + query
	c.Request = httptest.NewRequest(http.MethodGet, target, nil)
	c.Params = gin.Params{{Key: "provider", Value: providerParam}}
	return c, recorder
}

func TestNewProviderController_ReusesProvidedInstance(t *testing.T) {
	existing := &ProviderController{}
	assert.Same(t, existing, NewProviderController(existing))
}

func TestProviderUrl_InvalidReturnUrl(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("google", "returnUrl=not-relative")

	ctl.ProviderUrl(c)

	assert.NotEqual(t, http.StatusFound, recorder.Code)
}

func TestProviderUrl_InvalidClaimsType(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("google", "")
	c.Set("user", "not-a-claims-pointer")

	ctl.ProviderUrl(c)

	assert.NotEqual(t, http.StatusFound, recorder.Code)
}

func TestProviderUrl_InvalidProvider(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("not-a-provider", "")

	ctl.ProviderUrl(c)

	// "invalid provider" is a plain errors.New(...), not a registered custom
	// error, so HandleJSONResponse maps it to 500 (swagger doc says 400).
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProviderUrl_Success_Redirects(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("google", "returnUrl=/dashboard")

	ctl.ProviderUrl(c)

	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Location"), "https://accounts.google.com")
}

func encryptedState(t *testing.T, userId, returnUrl string) string {
	t.Helper()
	raw, err := json.Marshal(map[string]string{"userId": userId, "returnUrl": returnUrl})
	require.NoError(t, err)
	encrypted, err := encryption.Encrypt(string(raw))
	require.NoError(t, err)
	return encrypted
}

func TestProviderCallback_MissingState(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("google", "code=abc")

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProviderCallback_UndecryptableState(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	c, recorder := newProviderTestContext("google", "code=abc&state=not-encrypted!!")

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProviderCallback_InvalidStateJSON(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	encrypted, err := encryption.Encrypt("not-json")
	require.NoError(t, err)
	c, recorder := newProviderTestContext("google", "code=abc&state="+url.QueryEscape(encrypted))

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProviderCallback_EmptyStateObject(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	encrypted, err := encryption.Encrypt("{}")
	require.NoError(t, err)
	c, recorder := newProviderTestContext("google", "code=abc&state="+url.QueryEscape(encrypted))

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProviderCallback_InvalidProvider_RedirectsWithError(t *testing.T) {
	ctl := &ProviderController{signinService: newTestProviderService(t)}
	state := encryptedState(t, "", "/dashboard")
	c, recorder := newProviderTestContext("not-a-provider", "code=abc&state="+url.QueryEscape(state))

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Location"), "error=")
}

func TestProviderCallback_Success_SetsCookiesAndRedirects(t *testing.T) {
	s := newTestProviderService(t)
	email := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "ctrluser", email)
	restore := setGoogleURLs(server.URL+"/token", server.URL+"/userinfo")
	defer restore()

	called := stubSMTPAwait(t)

	ctl := &ProviderController{signinService: s}
	state := encryptedState(t, "", "/dashboard")
	c, recorder := newProviderTestContext("google", "code=good-code&state="+url.QueryEscape(state))

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
	awaitSMTP(t, called)
}

func TestProviderCallback_LinkExistingAccount(t *testing.T) {
	s := newTestProviderService(t)
	loggedInEmail := uniqueEmail(t)
	var loggedInAccount model.Account
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &loggedInEmail}, &loggedInAccount))

	newProviderEmail := uniqueEmail(t)
	server := newOAuthSuccessServer(t, "google-sub-"+uuid.NewString(), "linkeduser", newProviderEmail)
	restore := setGoogleURLs(server.URL+"/token", server.URL+"/userinfo")
	defer restore()

	ctl := &ProviderController{signinService: s}
	state := encryptedState(t, loggedInAccount.Id.String(), "/settings")
	c, recorder := newProviderTestContext("google", "code=good-code&state="+url.QueryEscape(state))

	ctl.ProviderCallback(c)

	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.NotContains(t, recorder.Header().Get("Location"), "error=")
}
