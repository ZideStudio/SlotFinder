package signin

type SigninDto struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type TokenResponseDto struct {
	AccessToken string `json:"access_token"`
}
