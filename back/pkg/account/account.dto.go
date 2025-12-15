package account

type AccountCreateDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountUpdateDto struct {
	UserName *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
