package helpers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Error   bool   `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ShouldBindJSON(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func ResponseJSON(c *gin.Context, obj any, err error) {
	if err != nil {
		code, message := parseError(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiError{
			Error:   true,
			Code:    code,
			Message: message,
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, obj)
}

func parseError(err error) (code, message string) {
	errorMessage := err.Error()
	errorParts := strings.Split(errorMessage, ": ")

	if len(errorParts) < 2 {
		code = "UNKNOWN_ERROR"
		message = errorMessage
		return
	}

	code = errorParts[0]
	message = errorParts[1]

	return
}
