package account

import "app/commons/constants"

// AccountProviderDto - linked OAuth provider
type AccountProviderDto struct {
	Provider constants.Provider `json:"provider"`
}

// AccountResponseDto - GET /account/me and PATCH /account
type AccountResponseDto struct {
	Username      *string                   `json:"username"`
	Email         *string                   `json:"email"`
	AvatarUrl     string                    `json:"avatarUrl"`
	Language      constants.AccountLanguage `json:"language"`
	Color         string                    `json:"color"`
	TimeZone      string                    `json:"timeZone"`
	Providers     []AccountProviderDto      `json:"providers"`
	TermsAccepted bool                      `json:"termsAccepted"`
	TermsVersion  *string                   `json:"termsVersion"`
}
