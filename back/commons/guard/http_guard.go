package guard

import (
	"app/commons/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxUploadSizeMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			limit,
		)

		if err := c.Request.ParseMultipartForm(limit); err != nil {
			c.Abort()
			helpers.HandleJSONResponse(c, nil, errors.New("file too large"))
			return
		}

		c.Next()
	}
}
