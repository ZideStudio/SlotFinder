package account

import (
	"app/commons/guard"
	model "app/db/models"
	"app/db/repository"
	"app/testutils"
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAccountTestContext(method string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}
	c.Request = httptest.NewRequest(method, "/", reader)
	c.Request.Header.Set("Content-Type", "application/json")
	return c, recorder
}

func newAuthenticatedAccountContext(t *testing.T, method string, body []byte, userId uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	c, recorder := newAccountTestContext(method, body)
	c.Set("user", &guard.Claims{Id: userId})
	return c, recorder
}

func TestAccountController_Create_InvalidBody(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodPost, []byte(`not-json`))

	ctl.Create(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountController_Create_ServiceError(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	// Passes body-binding validation but fails the service's password check, unlike InvalidBody above.
	dto := validCreateDto(t)
	dto.Password = "short"
	body, _ := json.Marshal(dto)
	c, recorder := newAccountTestContext(http.MethodPost, body)

	ctl.Create(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountController_Create_Success(t *testing.T) {
	s := newTestAccountService(t)
	called := stubSMTPAwait(t, s.mailService)
	ctl := &AccountController{accountService: s}
	dto := validCreateDto(t)
	body, _ := json.Marshal(dto)
	c, recorder := newAccountTestContext(http.MethodPost, body)

	ctl.Create(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
	awaitSMTP(t, called)
}

func TestAccountController_GetMe_InvalidClaimsType(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodGet, nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.GetMe(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_GetMe_Unauthenticated(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodGet, nil)

	ctl.GetMe(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_GetMe_Success(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "ctrl-getme-"+uuid.NewString())
	ctl := &AccountController{accountService: s}
	c, recorder := newAuthenticatedAccountContext(t, http.MethodGet, nil, account.Id)

	ctl.GetMe(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAccountController_Update_InvalidClaimsType(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodPatch, []byte(`{}`))
	c.Set("user", "not-a-claims-pointer")

	ctl.Update(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_Update_Unauthenticated(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodPatch, []byte(`{}`))

	ctl.Update(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_Update_InvalidBody(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "ctrl-updbad-"+uuid.NewString())
	ctl := &AccountController{accountService: s}
	c, recorder := newAuthenticatedAccountContext(t, http.MethodPatch, []byte(`not-json`), account.Id)

	ctl.Update(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountController_Update_Success_SetsCookies(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "ctrl-updok-"+uuid.NewString())
	ctl := &AccountController{accountService: s}

	newName := "upd-" + uuid.NewString()[:8]
	body, _ := json.Marshal(AccountUpdateDto{Username: &newName})
	c, recorder := newAuthenticatedAccountContext(t, http.MethodPatch, body, account.Id)

	ctl.Update(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}

func TestAccountController_Update_NoTokenRegen_NoCookies(t *testing.T) {
	s := newTestAccountService(t)
	account := createTestAccountWithUsername(t, s, "ctrl-updcolor-"+uuid.NewString())
	ctl := &AccountController{accountService: s}

	color := "#654321"
	body, _ := json.Marshal(AccountUpdateDto{Color: &color})
	c, recorder := newAuthenticatedAccountContext(t, http.MethodPatch, body, account.Id)

	ctl.Update(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Set-Cookie"))
}

func multipartImageRequest(t *testing.T, fieldName, fileName string, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.FormDataContentType()
}

func validPNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			img.Set(x, y, color.RGBA{R: 200, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func newUploadAvatarContext(t *testing.T, userId uuid.UUID, body *bytes.Buffer, contentType string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPatch, "/", body)
	c.Request.Header.Set("Content-Type", contentType)
	if userId != uuid.Nil {
		c.Set("user", &guard.Claims{Id: userId})
	}
	return c, recorder
}

func TestAccountController_UploadAvatar_InvalidClaimsType(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	c, recorder := newUploadAvatarContext(t, uuid.Nil, &bytes.Buffer{}, "multipart/form-data")
	c.Set("user", "not-a-claims-pointer")

	ctl.UploadAvatar(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_UploadAvatar_Unauthenticated(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	c, recorder := newUploadAvatarContext(t, uuid.Nil, &bytes.Buffer{}, "multipart/form-data")

	ctl.UploadAvatar(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_UploadAvatar_MissingFile(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	body, contentType := multipartImageRequest(t, "not-image", "avatar.png", []byte("data"))
	c, recorder := newUploadAvatarContext(t, uuid.New(), body, contentType)

	ctl.UploadAvatar(c)

	// "missing image" is a plain errors.New(...), not one of the sentinel
	// errors in constants.CUSTOM_ERRORS_MAP, so HandleJSONResponse maps it
	// to ERR_SERVER_ERROR/500 rather than the 400 the swagger doc implies.
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_UploadAvatar_InvalidImage(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	body, contentType := multipartImageRequest(t, "image", "avatar.png", []byte("not-a-real-image"))
	c, recorder := newUploadAvatarContext(t, uuid.New(), body, contentType)

	ctl.UploadAvatar(c)

	// Same as above: "invalid image file" isn't a registered custom error.
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_UploadAvatar_UnsupportedFormat(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}

	// GIF decodes fine (this file imports image/gif) but isn't in ALLOWED_PICTURE_FORMATS.
	img := image.NewPaletted(image.Rect(0, 0, 4, 4), []color.Color{color.White, color.Black})
	var buf bytes.Buffer
	require.NoError(t, gif.Encode(&buf, img, nil))

	body, contentType := multipartImageRequest(t, "image", "avatar.gif", buf.Bytes())
	c, recorder := newUploadAvatarContext(t, uuid.New(), body, contentType)

	ctl.UploadAvatar(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_UploadAvatar_Success(t *testing.T) {
	accountRepo := repository.NewAccountRepository(nil)
	account := model.Account{Id: uuid.New()}
	require.NoError(t, accountRepo.Create(repository.AccountCreateDto{Id: account.Id}, &model.Account{}))

	ctl := &AccountController{avatarService: &AvatarService{accountRepository: accountRepo}}
	body, contentType := multipartImageRequest(t, "image", "avatar.png", validPNGBytes(t))
	c, recorder := newUploadAvatarContext(t, account.Id, body, contentType)

	ctl.UploadAvatar(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var found model.Account
	require.NoError(t, accountRepo.FindOneById(account.Id, &found))
	assert.NotEmpty(t, found.AvatarData)
}

func TestAccountController_UploadAvatar_RepositoryError(t *testing.T) {
	ctl := &AccountController{avatarService: &AvatarService{accountRepository: repository.NewAccountRepository(testutils.ClosedDB(t))}}
	body, contentType := multipartImageRequest(t, "image", "avatar.png", validPNGBytes(t))
	c, recorder := newUploadAvatarContext(t, uuid.New(), body, contentType)

	ctl.UploadAvatar(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestAccountController_ForgotPassword_InvalidBody(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodPost, []byte(`not-json`))

	ctl.ForgotPassword(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountController_ForgotPassword_Success(t *testing.T) {
	s := newTestAccountService(t)
	stubSMTP(t, s.mailService)
	email := uniqueEmail(t)
	require.NoError(t, s.accountRepository.Create(repository.AccountCreateDto{Id: uuid.New(), Email: &email}, &model.Account{}))
	ctl := &AccountController{accountService: s}

	body, _ := json.Marshal(ForgotPasswordDto{Email: email})
	c, recorder := newAccountTestContext(http.MethodPost, body)

	ctl.ForgotPassword(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAccountController_ResetPassword_InvalidBody(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	c, recorder := newAccountTestContext(http.MethodPost, []byte(`not-json`))

	ctl.ResetPassword(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountController_ResetPassword_ServiceError(t *testing.T) {
	ctl := &AccountController{accountService: newTestAccountService(t)}
	body, _ := json.Marshal(ResetPasswordDto{Token: "bad-token", Password: "short"})
	c, recorder := newAccountTestContext(http.MethodPost, body)

	ctl.ResetPassword(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestNewAccountController_ReusesProvidedInstance(t *testing.T) {
	existing := &AccountController{}
	assert.Same(t, existing, NewAccountController(existing))
}

func TestNewAccountController_Nil_BuildsDefault(t *testing.T) {
	ctl := NewAccountController(nil)
	assert.NotNil(t, ctl.accountService)
	assert.NotNil(t, ctl.avatarService)
}

func newGetAvatarContext(accountId string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: "accountId", Value: accountId}}
	return c, recorder
}

func TestAccountController_GetAvatar_InvalidId(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	c, _ := newGetAvatarContext("not-a-uuid")

	ctl.GetAvatar(c)

	// GetAvatar's error paths call bare c.Status(...) with no body, so gin
	// only tracks the status internally (c.Writer.Status()) without flushing
	// it to the recorder — recorder.Code would stay at its 200 default.
	assert.Equal(t, http.StatusNotFound, c.Writer.Status())
}

func TestAccountController_GetAvatar_NotFound(t *testing.T) {
	ctl := &AccountController{avatarService: NewAvatarService(nil)}
	c, _ := newGetAvatarContext(uuid.New().String())

	ctl.GetAvatar(c)

	assert.Equal(t, http.StatusNotFound, c.Writer.Status())
}

func TestAccountController_GetAvatar_RepositoryError(t *testing.T) {
	ctl := &AccountController{avatarService: &AvatarService{accountRepository: repository.NewAccountRepository(testutils.ClosedDB(t))}}
	c, _ := newGetAvatarContext(uuid.New().String())

	ctl.GetAvatar(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}

func TestAccountController_GetAvatar_Success(t *testing.T) {
	accountRepo := repository.NewAccountRepository(nil)
	account := model.Account{Id: uuid.New(), AvatarData: []byte("avatar-bytes")}
	require.NoError(t, accountRepo.Create(repository.AccountCreateDto{Id: account.Id}, &model.Account{}))
	require.NoError(t, accountRepo.Updates(model.Account{Id: account.Id, AvatarData: account.AvatarData}))

	ctl := &AccountController{avatarService: &AvatarService{accountRepository: accountRepo}}
	c, recorder := newGetAvatarContext(account.Id.String())

	ctl.GetAvatar(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "avatar-bytes", recorder.Body.String())
	assert.NotEmpty(t, recorder.Header().Get("Last-Modified"))
	assert.Equal(t, "public, max-age=86400", recorder.Header().Get("Cache-Control"))
}

func TestAccountController_GetAvatar_NotModified(t *testing.T) {
	accountRepo := repository.NewAccountRepository(nil)
	account := model.Account{Id: uuid.New()}
	require.NoError(t, accountRepo.Create(repository.AccountCreateDto{Id: account.Id}, &model.Account{}))
	require.NoError(t, accountRepo.Updates(model.Account{Id: account.Id, AvatarData: []byte("avatar-bytes")}))

	ctl := &AccountController{avatarService: &AvatarService{accountRepository: accountRepo}}
	c, _ := newGetAvatarContext(account.Id.String())
	c.Request.Header.Set("If-Modified-Since", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))

	ctl.GetAvatar(c)

	assert.Equal(t, http.StatusNotModified, c.Writer.Status())
}
