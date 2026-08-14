package guard

import (
	"app/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTinyRSAKey builds a valid RSA key well under 1024 bits by hand,
// since rsa.GenerateKey refuses to produce keys that small.
func generateTinyRSAKey(t *testing.T, bits int) *rsa.PrivateKey {
	t.Helper()
	for {
		p, err := rand.Prime(rand.Reader, bits/2)
		require.NoError(t, err)
		q, err := rand.Prime(rand.Reader, bits/2)
		require.NoError(t, err)
		if p.Cmp(q) == 0 {
			continue
		}

		n := new(big.Int).Mul(p, q)
		phi := new(big.Int).Mul(
			new(big.Int).Sub(p, big.NewInt(1)),
			new(big.Int).Sub(q, big.NewInt(1)),
		)
		e := big.NewInt(65537)
		d := new(big.Int).ModInverse(e, phi)
		if d == nil {
			continue // e and phi weren't coprime — retry with fresh primes
		}

		key := &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{N: n, E: int(e.Int64())},
			D:         d,
			Primes:    []*big.Int{p, q},
		}
		key.Precompute()
		return key
	}
}

func newGuardTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, recorder
}

func TestGetUserClaims_NotSet(t *testing.T) {
	c, _ := newGuardTestContext()
	var user *Claims
	err := GetUserClaims(c, &user)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestGetUserClaims_ValidClaims(t *testing.T) {
	c, _ := newGuardTestContext()
	expected := &Claims{Id: uuid.New()}
	c.Set("user", expected)

	var user *Claims
	err := GetUserClaims(c, &user)
	assert.NoError(t, err)
	assert.Same(t, expected, user)
}

func TestGetUserClaims_WrongType(t *testing.T) {
	c, _ := newGuardTestContext()
	c.Set("user", "not-a-claims-pointer")

	var user *Claims
	err := GetUserClaims(c, &user)
	assert.Error(t, err)
}

func validClaims() *Claims {
	return &Claims{Id: uuid.New()}
}

func TestGenerateAccessToken_Success(t *testing.T) {
	token, err := GenerateAccessToken(validClaims())
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken_Success(t *testing.T) {
	token, err := GenerateAccessToken(validClaims())
	require.NoError(t, err)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("not-a-real-jwt")
	assert.Error(t, err)
}

// signWithTestKey signs claims with the repo's test RSA private key,
// bypassing GenerateAccessToken so the caller can set a custom ExpiresAt.
func signWithTestKey(t *testing.T, claims *Claims) string {
	t.Helper()
	c := config.GetConfig()
	privateKeyFile, err := os.ReadFile(c.Auth.PrivatePemPath)
	require.NoError(t, err)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return signed
}

func TestParseToken_ExpiredToken(t *testing.T) {
	claims := validClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	signed := signWithTestKey(t, claims)

	_, err := ParseToken(signed)
	assert.Error(t, err)
}

func TestAuthCheck_NoCookie_RequireAuth(t *testing.T) {
	c, recorder := newGuardTestContext()
	AuthCheck(nil)(c)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthCheck_NoCookie_AuthNotRequired(t *testing.T) {
	c, recorder := newGuardTestContext()
	AuthCheck(&AuthCheckParams{RequireAuthentication: false})(c)
	assert.NotEqual(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthCheck_InvalidToken(t *testing.T) {
	c, recorder := newGuardTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "not-a-real-jwt"})

	AuthCheck(nil)(c)
	assert.NotEqual(t, http.StatusOK, recorder.Code)
}

func TestAuthCheck_MissingUsername(t *testing.T) {
	token, err := GenerateAccessToken(validClaims())
	require.NoError(t, err)

	c, recorder := newGuardTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	AuthCheck(nil)(c)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAuthCheck_TermsNotAccepted(t *testing.T) {
	username := "someone"
	token, err := GenerateAccessToken(&Claims{Id: uuid.New(), Username: &username})
	require.NoError(t, err)

	c, recorder := newGuardTestContext()
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	AuthCheck(nil)(c)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAuthCheck_Success(t *testing.T) {
	username := "someone"
	token, err := GenerateAccessToken(&Claims{Id: uuid.New(), Username: &username, TermsAccepted: true})
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	nextCalled := false
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthCheck(nil))
	router.GET("/", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })
	router.ServeHTTP(recorder, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthCheck_SkipsCompleteProfileCheck(t *testing.T) {
	token, err := GenerateAccessToken(validClaims())
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthCheck(&AuthCheckParams{RequireAuthentication: true, RequireCompleteProfile: false}))
	router.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestParseToken_PublicKeyFileMissing(t *testing.T) {
	cfg := config.GetConfig()
	original := cfg.Auth.PublicPemPath
	cfg.Auth.PublicPemPath = "/nonexistent/public.pem"
	defer func() { cfg.Auth.PublicPemPath = original }()

	_, err := ParseToken("any-token")
	assert.Error(t, err)
}

func TestGenerateAccessToken_PrivateKeyFileMissing(t *testing.T) {
	cfg := config.GetConfig()
	original := cfg.Auth.PrivatePemPath
	cfg.Auth.PrivatePemPath = "/nonexistent/private.pem"
	defer func() { cfg.Auth.PrivatePemPath = original }()

	_, err := GenerateAccessToken(validClaims())
	assert.Error(t, err)
}

func TestGenerateAccessToken_PrivateKeyFileInvalid(t *testing.T) {
	cfg := config.GetConfig()
	original := cfg.Auth.PrivatePemPath

	invalidPath := filepath.Join(t.TempDir(), "invalid.pem")
	require.NoError(t, os.WriteFile(invalidPath, []byte("not a pem file"), 0o600))
	cfg.Auth.PrivatePemPath = invalidPath
	defer func() { cfg.Auth.PrivatePemPath = original }()

	_, err := GenerateAccessToken(validClaims())
	assert.Error(t, err)
}

func TestGenerateAccessToken_PrivateKeyTooSmallToSign(t *testing.T) {
	cfg := config.GetConfig()
	original := cfg.Auth.PrivatePemPath

	tinyKey := generateTinyRSAKey(t, 384)
	tinyKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(tinyKey),
	})
	tinyKeyPath := filepath.Join(t.TempDir(), "tiny.pem")
	require.NoError(t, os.WriteFile(tinyKeyPath, tinyKeyPem, 0o600))

	cfg.Auth.PrivatePemPath = tinyKeyPath
	defer func() { cfg.Auth.PrivatePemPath = original }()

	_, err := GenerateAccessToken(validClaims())
	assert.Error(t, err)
}

func TestShouldRenewToken_NilExpiresAt(t *testing.T) {
	assert.False(t, ShouldRenewToken(&Claims{}))
}

func TestAuthCheck_RenewsTokenNearExpiry(t *testing.T) {
	username := "someone"
	claims := &Claims{Id: uuid.New(), Username: &username, TermsAccepted: true}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(2 * time.Minute))
	signed := signWithTestKey(t, claims)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: signed})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthCheck(nil))
	router.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Set-Cookie"), "access_token=")
}
