package testutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileExists_ExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "some-file")
	require.NoError(t, os.WriteFile(path, []byte("data"), 0o600))

	assert.True(t, fileExists(path))
}

func TestFileExists_MissingPath(t *testing.T) {
	assert.False(t, fileExists(filepath.Join(t.TempDir(), "does-not-exist")))
}

func TestFileExists_Directory(t *testing.T) {
	assert.False(t, fileExists(t.TempDir()))
}

func readRSAKeyPair(t *testing.T, privateKeyPath, publicKeyPath string) {
	t.Helper()

	privateKeyPem, err := os.ReadFile(privateKeyPath)
	require.NoError(t, err)
	privateBlock, _ := pem.Decode(privateKeyPem)
	require.NotNil(t, privateBlock, "private key file should contain a valid PEM block")
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	require.NoError(t, err)

	publicKeyPem, err := os.ReadFile(publicKeyPath)
	require.NoError(t, err)
	publicBlock, _ := pem.Decode(publicKeyPem)
	require.NotNil(t, publicBlock, "public key file should contain a valid PEM block")
	publicKeyAny, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	require.NoError(t, err)
	publicKey, ok := publicKeyAny.(*rsa.PublicKey)
	require.True(t, ok, "expected an RSA public key")

	assert.True(t, privateKey.PublicKey.Equal(publicKey), "public key should correspond to the generated private key")
}

func TestEnsureTestJWTKeyPair_GeneratesValidKeyPair(t *testing.T) {
	dir := t.TempDir()
	privateKeyPath := filepath.Join(dir, "private.pem")
	publicKeyPath := filepath.Join(dir, "public.pem")

	require.NoError(t, EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath))

	assert.True(t, fileExists(privateKeyPath))
	assert.True(t, fileExists(publicKeyPath))
	readRSAKeyPair(t, privateKeyPath, publicKeyPath)
}

func TestEnsureTestJWTKeyPair_AlreadyExists_NoOp(t *testing.T) {
	dir := t.TempDir()
	privateKeyPath := filepath.Join(dir, "private.pem")
	publicKeyPath := filepath.Join(dir, "public.pem")

	require.NoError(t, EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath))
	privateBefore, err := os.ReadFile(privateKeyPath)
	require.NoError(t, err)
	publicBefore, err := os.ReadFile(publicKeyPath)
	require.NoError(t, err)

	require.NoError(t, EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath))

	privateAfter, err := os.ReadFile(privateKeyPath)
	require.NoError(t, err)
	publicAfter, err := os.ReadFile(publicKeyPath)
	require.NoError(t, err)

	assert.Equal(t, privateBefore, privateAfter, "an existing key pair should not be regenerated")
	assert.Equal(t, publicBefore, publicAfter, "an existing key pair should not be regenerated")
}

func TestEnsureTestJWTKeyPair_MkdirAllFails(t *testing.T) {
	dir := t.TempDir()
	// A file where a directory component is expected makes os.MkdirAll fail.
	blocker := filepath.Join(dir, "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("not a directory"), 0o600))

	privateKeyPath := filepath.Join(blocker, "subdir", "private.pem")
	publicKeyPath := filepath.Join(blocker, "subdir", "public.pem")

	err := EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath)
	assert.Error(t, err)
}

func TestEnsureTestJWTKeyPair_WriteFileFails(t *testing.T) {
	dir := t.TempDir()
	// A directory at the target path makes os.WriteFile fail.
	privateKeyPath := filepath.Join(dir, "private-as-dir")
	require.NoError(t, os.MkdirAll(privateKeyPath, 0o755))
	publicKeyPath := filepath.Join(dir, "public.pem")

	err := EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath)
	assert.Error(t, err)
}

func TestEnsureTestJWTKeyPair_OnlyPrivateKeyPresent_RegeneratesPair(t *testing.T) {
	dir := t.TempDir()
	privateKeyPath := filepath.Join(dir, "private.pem")
	publicKeyPath := filepath.Join(dir, "public.pem")

	require.NoError(t, EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath))
	privateBefore, err := os.ReadFile(privateKeyPath)
	require.NoError(t, err)
	require.NoError(t, os.Remove(publicKeyPath))

	require.NoError(t, EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath))

	assert.True(t, fileExists(publicKeyPath), "the missing public key should be regenerated")
	privateAfter, err := os.ReadFile(privateKeyPath)
	require.NoError(t, err)
	assert.NotEqual(t, privateBefore, privateAfter, "a fresh pair should be generated when only one of the two files is present")
	readRSAKeyPair(t, privateKeyPath, publicKeyPath)
}
