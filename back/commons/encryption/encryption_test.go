package encryption

import (
	"os"
	"testing"
)

const testKey = "1234567890abcdef"

func TestEncryptDecrypt(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	plainText := "Secret message"
	encrypted, err := Encrypt(plainText)
	if err != nil {
		t.Fatalf("Error during encryption: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypted text must not be empty")
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Error during decryption: %v", err)
	}

	if decrypted != plainText {
		t.Errorf("Decrypted text does not match the original.\nExpected: %q\nGot: %q", plainText, decrypted)
	}
}

func TestEncryptThenDecryptOriginalValue(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	original := "Original secret value"
	encrypted, err := Encrypt(original)
	if err != nil {
		t.Fatalf("Error during encryption: %v", err)
	}
	if encrypted == "" {
		t.Fatal("Encrypted text must not be empty")
	}
	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Error during decryption: %v", err)
	}
	if decrypted != original {
		t.Errorf("Decrypted value does not match the original.\nExpected: %q\nGot: %q", original, decrypted)
	}
}

func TestDecryptWithInvalidBase64(t *testing.T) {
	_ = os.Setenv("ENCRYPTION_KEY", testKey)

	_, err := Decrypt("!!not_base64!!")
	if err == nil {
		t.Error("Decryption of non-base64 string should fail")
	}
}
