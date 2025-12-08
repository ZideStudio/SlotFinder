package account

import "app/commons/constants"

type AccountCreateDto struct {
	Email         string                    `json:"email" binding:"required,email,min=6,max=320"`
	Password      string                    `json:"password" binding:"required,max=100"`
	Language      constants.AccountLanguage `json:"language" binding:"required,oneof=en fr"`
	TermsAccepted bool                      `json:"termsAccepted" binding:"required,eq=true"`
	TermsVersion  string                    `json:"termsVersion" binding:"required"`
}

type AccountUpdateDto struct {
	UserName      *string                    `json:"username" binding:"omitempty,min=3,max=30"`
	Email         *string                    `json:"email" binding:"omitempty,min=6,max=320"`
	Password      *string                    `json:"password" binding:"omitempty,max=100"`
	Color         *string                    `json:"color"`
	Language      *constants.AccountLanguage `json:"language" binding:"omitempty,oneof=en fr"`
	TermsAccepted *bool                      `json:"termsAccepted" binding:"omitempty,eq=true"`
	TermsVersion  *string                    `json:"termsVersion"`
}

type ForgotPasswordDto struct {
	Email string `json:"email" binding:"required,email,min=6,max=320"`
}

type ResetPasswordDto struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,max=100"`
}

type AccountTokensDto struct {
	AccessToken  string
	RefreshToken string
}
