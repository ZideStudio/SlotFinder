package encryption

import (
	"crypto/rand"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKey = "1234567890abcdef"

type failingReader struct{}

func (failingReader) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestEncryptDecrypt(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	plainText := "Secret message"
	encrypted, err := Encrypt(plainText)
	assert.NoError(t, err, "Error during encryption")
	assert.NotEmpty(t, encrypted, "Encrypted text must not be empty")

	decrypted, err := Decrypt(encrypted)
	assert.NoError(t, err, "Error during decryption")
	assert.Equal(t, plainText, decrypted, "Decrypted text does not match the original")
}

func TestDecryptWithInvalidBase64(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	_, err := Decrypt("!!not_base64!!")
	if assert.Error(t, err, "Decryption of non-base64 string should fail") {
		assert.Contains(t, err.Error(), "illegal base64", "Error message should mention base64 decoding")
	}
}

func TestEncrypt_InvalidKey(t *testing.T) {
	original := os.Getenv("ENCRYPTION_KEY")
	defer os.Setenv("ENCRYPTION_KEY", original)
	_ = os.Setenv("ENCRYPTION_KEY", "too-short")

	_, err := Encrypt("secret")
	assert.ErrorContains(t, err, "ENCRYPTION_KEY must be")
}

func TestDecrypt_InvalidKey(t *testing.T) {
	original := os.Getenv("ENCRYPTION_KEY")
	defer os.Setenv("ENCRYPTION_KEY", original)
	_ = os.Setenv("ENCRYPTION_KEY", "too-short")

	_, err := Decrypt("anything")
	assert.ErrorContains(t, err, "ENCRYPTION_KEY must be")
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	_, err := Decrypt("c2hvcnQ=") // base64("short"), shorter than aes.BlockSize
	assert.ErrorContains(t, err, "ciphertext too short")
}

func TestEncrypt_RandReaderFails(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	original := rand.Reader
	rand.Reader = failingReader{}
	defer func() { rand.Reader = original }()

	_, err := Encrypt("secret")
	assert.ErrorContains(t, err, "read failed")
}
