package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAccount_ComparePassword_NilPassword(t *testing.T) {
	account := Account{Password: nil}
	assert.False(t, account.ComparePassword("anything"))
}

func TestAccount_ComparePassword_CorrectPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	require.NoError(t, err)
	hashed := string(hash)
	account := Account{Password: &hashed}

	assert.True(t, account.ComparePassword("correct-password"))
}

func TestAccount_ComparePassword_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	require.NoError(t, err)
	hashed := string(hash)
	account := Account{Password: &hashed}

	assert.False(t, account.ComparePassword("wrong-password"))
}

func TestAccount_Sanitized_WithOverrideColor(t *testing.T) {
	username := "alice"
	overrideColor := "#ffffff"
	account := Account{UserName: &username, Color: "#000000"}

	sanitized := account.Sanitized(&overrideColor)

	assert.Equal(t, "#ffffff", sanitized.Color)
}

func TestAccount_TableName(t *testing.T) {
	assert.Equal(t, "account", Account{}.TableName())
}

func TestAccountEvent_TableName(t *testing.T) {
	assert.Equal(t, "account_event", AccountEvent{}.TableName())
}

func TestAccountEvent_Sanitized(t *testing.T) {
	username := "alice"
	ownerName := "bob"
	color := "#123456"
	accountEvent := AccountEvent{
		Color:   &color,
		Account: Account{UserName: &username},
		Event:   Event{Owner: Account{UserName: &ownerName}},
	}

	sanitized := accountEvent.Sanitized()

	assert.Equal(t, "alice", *sanitized.Account.UserName)
	assert.Equal(t, "#123456", sanitized.Account.Color)
	assert.Equal(t, "bob", *sanitized.Event.Owner.UserName)
}

func TestAccountProvider_TableName(t *testing.T) {
	assert.Equal(t, "account_provider", AccountProvider{}.TableName())
}

func TestAvailability_TableName(t *testing.T) {
	assert.Equal(t, "availability", Availability{}.TableName())
}

func TestRefreshToken_TableName(t *testing.T) {
	assert.Equal(t, "refresh_token", RefreshToken{}.TableName())
}
