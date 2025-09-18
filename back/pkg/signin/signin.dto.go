package signin

type SigninDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponseDto struct {
	AccessToken string `json:"access_token"`
}
