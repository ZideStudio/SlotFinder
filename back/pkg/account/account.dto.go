package account

type AccountCreateDto struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountUpdateDto struct {
	UserName *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Color    *string `json:"color"`
}
