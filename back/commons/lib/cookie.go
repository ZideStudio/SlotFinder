package lib

import (
	"app/commons/constants"
	"app/config"

	"github.com/gin-gonic/gin"
)

func setCookie(c *gin.Context, token string, expiration int) {
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

func AddCookie(c *gin.Context, token string) {
	setCookie(c, token, constants.TOKEN_EXPIRATION)
}

func RemoveCookie(c *gin.Context) {
	setCookie(c, "", -1)
}
