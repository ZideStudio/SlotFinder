package account

import model "app/db/models"

func MapToAccountResponseDto(a model.Account) AccountResponseDto {
	providers := make([]AccountProviderDto, 0, len(a.Providers))
	for _, p := range a.Providers {
		providers = append(providers, AccountProviderDto{Provider: p.Provider})
	}
	return AccountResponseDto{
		UserName:  a.UserName,
		Email:     a.Email,
		AvatarUrl: a.AvatarUrl,
		Language:  a.Language,
		Color:     a.Color,
		TimeZone:  a.TimeZone,
		Providers: providers,
	}
}
