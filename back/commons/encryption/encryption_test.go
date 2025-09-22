package encryption

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKey = "1234567890abcdef"

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
	assert.Error(t, err, "Decryption of non-base64 string should fail")
}
