package account

import "app/commons/constants"

type AccountCreateDto struct {
	Email    string                    `json:"email" binding:"required,email"`
	Password string                    `json:"password" binding:"required"`
	Language constants.AccountLanguage `json:"language" binding:"required,oneof=en fr"`
}

type AccountUpdateDto struct {
	UserName *string                    `json:"username"`
	Email    *string                    `json:"email"`
	Password *string                    `json:"password"`
	Color    *string                    `json:"color"`
	Language *constants.AccountLanguage `json:"language" binding:"omitempty,oneof=en fr"`
}

type ForgotPasswordDto struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordDto struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
