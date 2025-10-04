package lib

import (
	"app/config"

	"github.com/gin-gonic/gin"
)

func RemoveCookie(c *gin.Context) {
	c.SetCookie(
		"access_token",            // name
		"",                        // value
		-1,                        // max age in seconds
		"/",                       // path
		config.GetConfig().Domain, // domain
		true,                      // secure
		true,                      // httpOnly
	)
}
