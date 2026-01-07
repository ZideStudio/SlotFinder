package lib

import (
	"app/commons/constants"
	"time"

	"github.com/gin-gonic/gin"
)

func SetAccessTokenCookie(c *gin.Context, token string, expiration int) {
	if expiration == 0 { // same as refresh token expiration to keep it saved as long as it can be renewed
		expiration = int(constants.REFRESH_TOKEN_EXPIRATION / time.Second)
	}

	c.SetCookie(
		"access_token", // name
		token,          // value
		expiration,     // max age in seconds
		"/api",         // path
		"",             // domain (empty = current domain)
		true,           // secure
		true,           // httpOnly
	)
}

func SetRefreshTokenCookie(c *gin.Context, token string, expiration int) {
	if expiration == 0 {
		expiration = int(constants.REFRESH_TOKEN_EXPIRATION / time.Second)
	}

	c.SetCookie(
		"refresh_token",        // name
		token,                  // value
		expiration,             // max age in seconds
		"/api/v1/auth/refresh", // path
		"",                     // domain (empty = current domain)
		true,                   // secure
		true,                   // httpOnly
	)
}
