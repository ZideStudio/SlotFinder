package account

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapToAccountResponseDto_WithProviders(t *testing.T) {
	t.Parallel()
	username := "mapped"
	account := model.Account{
		Username: &username,
		Providers: []model.AccountProvider{
			{Provider: constants.PROVIDER_GOOGLE},
			{Provider: constants.PROVIDER_GITHUB},
		},
	}

	dto := MapToAccountResponseDto(account)

	assert.Equal(t, username, *dto.Username)
	assert.Len(t, dto.Providers, 2)
	assert.Equal(t, constants.PROVIDER_GOOGLE, dto.Providers[0].Provider)
	assert.Equal(t, constants.PROVIDER_GITHUB, dto.Providers[1].Provider)
}

func TestMapToAccountResponseDto_NoProviders(t *testing.T) {
	t.Parallel()
	dto := MapToAccountResponseDto(model.Account{})
	assert.Empty(t, dto.Providers)
}

func TestMapToAccountResponseDto_TermsAccepted(t *testing.T) {
	t.Parallel()
	acceptedAt := time.Now()
	version := "1.2"
	account := model.Account{
		TermsAcceptedAt: &acceptedAt,
		TermsVersion:    &version,
	}

	dto := MapToAccountResponseDto(account)

	assert.True(t, dto.TermsAccepted)
	require.NotNil(t, dto.TermsVersion)
	assert.Equal(t, version, *dto.TermsVersion)
}

func TestMapToAccountResponseDto_TermsNotAccepted(t *testing.T) {
	t.Parallel()
	dto := MapToAccountResponseDto(model.Account{})

	assert.False(t, dto.TermsAccepted)
	assert.Nil(t, dto.TermsVersion)
}
