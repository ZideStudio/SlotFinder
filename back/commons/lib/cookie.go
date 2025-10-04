package lib

import (
	"app/commons/constants"
	"app/config"

	"github.com/gin-gonic/gin"
)

func SetAccessTokenCookie(c *gin.Context, token string, expiration int) {
	if expiration == 0 {
		expiration = constants.TOKEN_EXPIRATION
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
