package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueEmail_ReturnsDistinctAddresses(t *testing.T) {
	assert.NotEqual(t, UniqueEmail(t), UniqueEmail(t))
}

func TestStubSMTP_ReturnsNilWithoutSendingMail(t *testing.T) {
	var target SendMailFunc
	StubSMTP(t, &target)

	assert.NoError(t, target("addr", nil, "from", []string{"to@example.com"}, []byte("msg")))
}

func TestStubSMTPAwait_SignalsOnInvocation(t *testing.T) {
	var target SendMailFunc
	called := StubSMTPAwait(t, &target)

	assert.NoError(t, target("addr", nil, "from", []string{"to@example.com"}, []byte("msg")))
	AwaitSMTP(t, called)
}

func TestAwaitSMTP_FailsWhenNeverCalled(t *testing.T) {
	msg := expectFatal(t, func() {
		AwaitSMTP(fakeT{}, make(chan struct{}))
	})

	assert.Contains(t, msg, "expected the async SendMail goroutine to run")
}
