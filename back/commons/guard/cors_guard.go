package guard

import (
	"app/config"
	"slices"

	"github.com/gin-gonic/gin"
)

type CorsGuard struct {
}

func (g *CorsGuard) CorsCheck() gin.HandlerFunc {
	config := config.GetConfig()

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if slices.Contains(config.Origins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Content-Disposition, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}
