package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

func EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath string) error {
	privateKeyExists := fileExists(privateKeyPath)
	publicKeyExists := fileExists(publicKeyPath)
	if privateKeyExists && publicKeyExists {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(publicKeyPath), 0o755); err != nil {
		return err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if privateKeyPem == nil {
		return errors.New("failed to encode private key pem")
	}

	publicKeyDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	})
	if publicKeyPem == nil {
		return errors.New("failed to encode public key pem")
	}

	if err := os.WriteFile(privateKeyPath, privateKeyPem, 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(publicKeyPath, publicKeyPem, 0o644); err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
