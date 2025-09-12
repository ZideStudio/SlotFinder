package health

import (
	"app/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Ready(c *gin.Context) {
	if !db.TestConnection() {
		c.String(http.StatusServiceUnavailable, "db not ready")
		return
	}

	c.String(http.StatusOK, "ready")
}

func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "work")
}
