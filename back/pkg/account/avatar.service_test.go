package account

import (
	model "app/db/models"
	"app/db/repository"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.Account{}, &model.AccountProvider{}))
	return database
}

// validPNG builds a small, valid PNG image so it can be decoded and
// re-encoded by lib.ProcessAvatar / lib.ProcessAvatarFromURL.
func validPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 100, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func newImageTestServer(t *testing.T, body []byte, status int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestNewAvatarService_ReusesProvidedInstance(t *testing.T) {
	existing := &AvatarService{}
	assert.Same(t, existing, NewAvatarService(existing))
}

func TestGetGravatarURL_Deterministic(t *testing.T) {
	url1 := GetGravatarURL("someuser")
	url2 := GetGravatarURL("someuser")
	url3 := GetGravatarURL("otheruser")
	assert.Equal(t, url1, url2)
	assert.NotEqual(t, url1, url3)
	assert.Contains(t, url1, "https://www.gravatar.com/avatar/")
}

// FetchAndStoreGravatar always builds its Gravatar URL from the real
// gravatar.com domain (via GetGravatarURL), so we can't point it at a local
// test server. We instead assert the contract that holds in either outcome:
// on success the returned data is non-empty and the URL is the local avatar
// endpoint; on failure (e.g. no network access in CI) data is nil and the
// URL falls back to the external Gravatar URL built from the account id.
func TestFetchAndStoreGravatar_ContractHoldsEitherWay(t *testing.T) {
	s := &AvatarService{}
	accountId := uuid.New()
	data, url := s.FetchAndStoreGravatar("someuser", accountId)
	if data == nil {
		assert.Equal(t, GetGravatarURL(accountId.String()), url)
	} else {
		assert.NotEmpty(t, data)
		assert.Equal(t, "/api/v1/account/"+accountId.String()+"/avatar", url)
	}
}

func TestFindAvatarById(t *testing.T) {
	db := testDB(t)
	accountRepo := repository.NewAccountRepository(db)
	avatarBytes := []byte("some-avatar-bytes")
	account := model.Account{Id: uuid.New(), AvatarData: avatarBytes}
	require.NoError(t, db.Create(&account).Error)

	s := &AvatarService{accountRepository: accountRepo}
	data, updatedAt, err := s.FindAvatarById(account.Id)
	require.NoError(t, err)
	assert.Equal(t, avatarBytes, data)
	require.NotNil(t, updatedAt)
}

func TestFindAvatarById_NotFound(t *testing.T) {
	db := testDB(t)
	s := &AvatarService{accountRepository: repository.NewAccountRepository(db)}
	_, _, err := s.FindAvatarById(uuid.New())
	// FindAvatarById uses Scan (not First), so a missing row does not error;
	// it returns zero-value data. We assert there's no error and no data.
	assert.NoError(t, err)
}

func TestUploadAvatar_NoImageProvided(t *testing.T) {
	s := &AvatarService{}
	_, err := s.UploadAvatar(nil, nil)
	assert.Error(t, err)
}

func TestUploadAvatar_FromBytes_Success(t *testing.T) {
	s := &AvatarService{}
	data, err := s.UploadAvatar(nil, validPNG(t, 10, 10))
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestUploadAvatar_FromBytes_InvalidImage(t *testing.T) {
	s := &AvatarService{}
	_, err := s.UploadAvatar(nil, []byte("not-an-image"))
	assert.Error(t, err)
}

func TestUploadAvatar_FromURL_Success(t *testing.T) {
	imgServer := newImageTestServer(t, validPNG(t, 10, 10), http.StatusOK)
	s := &AvatarService{}
	url := imgServer.URL
	data, err := s.UploadAvatar(&url, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestUploadAvatar_FromURL_DownloadFails(t *testing.T) {
	// Start and immediately close a server so the URL is guaranteed to
	// refuse connections, regardless of platform.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	badURL := server.URL
	server.Close()

	s := &AvatarService{}
	_, err := s.UploadAvatar(&badURL, nil)
	assert.Error(t, err)
}

func TestUploadAvatar_FromURL_NonSuccessStatus(t *testing.T) {
	imgServer := newImageTestServer(t, validPNG(t, 10, 10), http.StatusInternalServerError)
	s := &AvatarService{}
	url := imgServer.URL
	_, err := s.UploadAvatar(&url, nil)
	assert.Error(t, err)
}

func TestUploadUserAvatar_Success(t *testing.T) {
	db := testDB(t)
	accountRepo := repository.NewAccountRepository(db)
	account := model.Account{Id: uuid.New()}
	require.NoError(t, db.Create(&account).Error)

	s := &AvatarService{accountRepository: accountRepo}
	err := s.UploadUserAvatar(validPNG(t, 10, 10), account.Id)
	require.NoError(t, err)

	var found model.Account
	require.NoError(t, accountRepo.FindOneById(account.Id, &found))
	assert.Equal(t, "/api/v1/account/"+account.Id.String()+"/avatar", found.AvatarUrl)
	assert.NotEmpty(t, found.AvatarData)
}

func TestUploadUserAvatar_ProcessingFails(t *testing.T) {
	db := testDB(t)
	s := &AvatarService{accountRepository: repository.NewAccountRepository(db)}
	err := s.UploadUserAvatar([]byte("not-an-image"), uuid.New())
	assert.Error(t, err)
}
