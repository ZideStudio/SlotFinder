package testutils

import (
	"fmt"
	"net/smtp"
	"testing"
	"time"

	"github.com/google/uuid"
)

// UniqueEmail returns a fresh, collision-free email address for tests that
// need to create an account without hardcoding an address other tests might
// also use.
func UniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s@example.com", uuid.NewString())
}

// SendMailFunc mirrors mail.MailService's field type, declared here to
// avoid an import cycle (pkg/mail's tests import testutils).
type SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// StubSMTP points *target (typically &m.SendMailFunc) at a no-op stub for
// the duration of the test, discarding whatever mail would have been sent.
func StubSMTP(t *testing.T, target *SendMailFunc) {
	t.Helper()
	*target = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error { return nil }
}

// StubSMTPAwait points *target (typically &m.SendMailFunc) at a stub and
// returns a channel that receives once the stub has actually been invoked,
// so callers can wait for an asynchronously-spawned SendMail goroutine.
func StubSMTPAwait(t *testing.T, target *SendMailFunc) <-chan struct{} {
	t.Helper()
	called := make(chan struct{}, 1)
	*target = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called <- struct{}{}
		return nil
	}
	return called
}

// AwaitSMTP blocks until called receives (see StubSMTPAwait) or fails the
// test after 2s, for callers waiting on an asynchronously-spawned SendMail
// goroutine.
func AwaitSMTP(t *testing.T, called <-chan struct{}) {
	t.Helper()
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("expected the async SendMail goroutine to run")
	}
}
