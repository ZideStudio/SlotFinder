package account

type AccountCreateDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountUpdateDto struct {
	UserName *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Color    *string `json:"color"`
}

type ForgotPasswordDto struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordDto struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
