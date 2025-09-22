package guard

import (
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type CorsGuard struct {
	allowedOrigins []string
}

func (g *CorsGuard) init() {
	if len(g.allowedOrigins) > 0 {
		return
	}

	originsString := os.Getenv("ORIGINS")
	var allowedOrigins []string
	if originsString != "" {
		allowedOrigins = strings.Split(originsString, ",")
	}

	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	g.allowedOrigins = allowedOrigins
}

func (g *CorsGuard) CorsCheck() gin.HandlerFunc {
	g.init()
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if slices.Contains(g.allowedOrigins, origin) {
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
