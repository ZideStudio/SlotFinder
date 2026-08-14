package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// EnsureTestJWTKeyPair is called from several packages' TestMain, each its
// own OS process run in parallel by `go test ./...`. lockPath (O_CREATE|O_EXCL)
// is a cross-process mutex so only one process generates the missing keypair.
func EnsureTestJWTKeyPair(privateKeyPath, publicKeyPath string) error {
	if fileExists(privateKeyPath) && fileExists(publicKeyPath) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0o755); err != nil {
		return err
	}

	lockPath := privateKeyPath + ".generating.lock"
	lockFile, err := acquireGenerationLock(lockPath, privateKeyPath, publicKeyPath)
	if err != nil {
		return err
	}
	if lockFile == nil {
		// Another process finished generating while we were waiting.
		return nil
	}
	defer func() {
		_ = lockFile.Close()
		_ = os.Remove(lockPath)
	}()

	// Re-check: another process could have generated the pair between our
	// first check above and acquiring the lock.
	if fileExists(privateKeyPath) && fileExists(publicKeyPath) {
		return nil
	}

	return generateTestJWTKeyPair(privateKeyPath, publicKeyPath)
}

// acquireGenerationLock returns a held lock file if the caller should
// generate the keypair, or (nil, nil) if another process already finished
// generating it while this one was waiting.
func acquireGenerationLock(lockPath, privateKeyPath, publicKeyPath string) (*os.File, error) {
	deadline := time.Now().Add(30 * time.Second)
	for {
		lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			return lockFile, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		if fileExists(privateKeyPath) && fileExists(publicKeyPath) {
			return nil, nil
		}
		if time.Now().After(deadline) {
			return nil, errors.New("timed out waiting for another process to generate the test JWT keypair")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func generateTestJWTKeyPair(privateKeyPath, publicKeyPath string) error {
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
