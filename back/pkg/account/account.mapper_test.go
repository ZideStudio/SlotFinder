package account

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"

	"github.com/stretchr/testify/assert"
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
