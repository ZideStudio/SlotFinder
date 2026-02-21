package signin

type SigninDto struct {
	Identifier string `json:"identifier" binding:"required,min=3,max=320"`
	Password   string `json:"password" binding:"required,max=100"`
}

type TokenResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
