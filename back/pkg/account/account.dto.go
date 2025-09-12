package account

import (
	model "app/db/models"
	"app/pkg/signin"
)

type AccountCreateDto struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountCreateResponseDto struct {
	model.Account
	signin.TokenResponseDto
}

type AccountUpdateDto struct {
	UserName *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
