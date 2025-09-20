package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func getKey() ([]byte, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	switch len(key) {
	case 16, 24, 32:
		return key, nil
	default:
		return nil, errors.New("ENCRYPTION_KEY must be 16, 24, or 32 bytes")
	}
}

func Encrypt(text string) (string, error) {
	key, err := getKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	ctr := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(text))
	ctr.XORKeyStream(cipherText, []byte(text))

	return base64.StdEncoding.EncodeToString(append(iv, cipherText...)), nil
}

func Decrypt(text string) (string, error) {
	key, err := getKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv, cipherText := data[:aes.BlockSize], data[aes.BlockSize:]
	ctr := cipher.NewCTR(block, iv)
	plainText := make([]byte, len(cipherText))
	ctr.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
