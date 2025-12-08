package lib

import (
	"app/commons/constants"
	"app/config"

	"github.com/gin-gonic/gin"
)

func SetAccessTokenCookie(c *gin.Context, token string, expiration int) {
	if expiration == 0 {
		expiration = constants.ACCESS_TOKEN_EXPIRATION
	}

	c.SetCookie(
		"access_token",            // name
		token,                     // value
		expiration,                // max age in seconds
		"/",                       // path
		config.GetConfig().Domain, // domain
		true,                      // secure
		true,                      // httpOnly
	)
}

func SetRefreshTokenCookie(c *gin.Context, token string, expiration int) {
	if expiration == 0 {
		expiration = constants.REFRESH_TOKEN_EXPIRATION
	}

	c.SetCookie(
		"refresh_token",           // name
		token,                     // value
		expiration,                // max age in seconds
		"/",                       // path
		config.GetConfig().Domain, // domain
		true,                      // secure
		true,                      // httpOnly
	)
}
