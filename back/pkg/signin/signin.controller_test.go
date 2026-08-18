package signin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSigninTestContext(method string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, recorder
}

func TestNewSigninController_ReusesProvidedInstance(t *testing.T) {
	existing := &SigninController{}
	assert.Same(t, existing, NewSigninController(existing))
}

func TestNewSigninController_Nil_BuildsRealDependencies(t *testing.T) {
	ctl := NewSigninController(nil)
	assert.NotNil(t, ctl.signinService)
}

func TestSigninController_Signin_InvalidBody(t *testing.T) {
	ctl := &SigninController{signinService: newTestSigninService(t)}
	c, recorder := newSigninTestContext(http.MethodPost, []byte(`not-json`))

	ctl.Signin(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSigninController_Signin_InvalidCredentials(t *testing.T) {
	ctl := &SigninController{signinService: newTestSigninService(t)}
	body, _ := json.Marshal(SigninDto{Identifier: "unknown-" + uuid.NewString(), Password: "whatever"})
	c, recorder := newSigninTestContext(http.MethodPost, body)

	ctl.Signin(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestSigninController_Signin_Success(t *testing.T) {
	s := newTestSigninService(t)
	username := "ctrlsignin-" + uuid.NewString()
	createTestAccount(t, s, username, "CorrectPassword123!")
	ctl := &SigninController{signinService: s}

	body, _ := json.Marshal(SigninDto{Identifier: username, Password: "CorrectPassword123!"})
	c, recorder := newSigninTestContext(http.MethodPost, body)

	ctl.Signin(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}

func TestSigninController_Refresh_NoCookie(t *testing.T) {
	ctl := &SigninController{signinService: newTestSigninService(t)}
	c, recorder := newSigninTestContext(http.MethodPost, nil)

	ctl.Refresh(c)

	assert.NotEqual(t, http.StatusOK, recorder.Code)
}

func TestSigninController_Refresh_InvalidToken_ClearsCookies(t *testing.T) {
	ctl := &SigninController{signinService: newTestSigninService(t)}
	c, recorder := newSigninTestContext(http.MethodPost, nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "not-a-real-token"})

	ctl.Refresh(c)

	assert.NotEqual(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}

func TestSigninController_Refresh_Success(t *testing.T) {
	s := newTestSigninService(t)
	username := "ctrlrefresh-" + uuid.NewString()
	accountId := createTestAccount(t, s, username, "")
	rawToken, err := s.refreshTokenRepository.Create(accountId, time.Now().Add(time.Hour))
	require.NoError(t, err)

	ctl := &SigninController{signinService: s}
	c, recorder := newSigninTestContext(http.MethodPost, nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: rawToken})

	ctl.Refresh(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}
