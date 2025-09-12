package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

var secret = os.Getenv("SECURED_STR_KEY")

// Encrypt
func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	ctr := cipher.NewCTR(block, bytes)
	cipherText := make([]byte, len(plainText))
	ctr.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Decrypt
func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	ctr := cipher.NewCTR(block, bytes)
	plainText := make([]byte, len(cipherText))
	ctr.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
