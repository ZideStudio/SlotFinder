package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_ComparePassword_NilPassword(t *testing.T) {
	account := Account{Password: nil}
	assert.False(t, account.ComparePassword("anything"))
}
