package account

import (
	"app/commons/constants"
	model "app/db/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToAccountResponseDto_WithProviders(t *testing.T) {
	username := "mapped"
	account := model.Account{
		UserName: &username,
		Providers: []model.AccountProvider{
			{Provider: constants.PROVIDER_GOOGLE},
			{Provider: constants.PROVIDER_GITHUB},
		},
	}

	dto := MapToAccountResponseDto(account)

	assert.Equal(t, username, *dto.UserName)
	assert.Len(t, dto.Providers, 2)
	assert.Equal(t, constants.PROVIDER_GOOGLE, dto.Providers[0].Provider)
	assert.Equal(t, constants.PROVIDER_GITHUB, dto.Providers[1].Provider)
}

func TestMapToAccountResponseDto_NoProviders(t *testing.T) {
	dto := MapToAccountResponseDto(model.Account{})
	assert.Empty(t, dto.Providers)
}
