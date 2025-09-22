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

func SetHttpContextBody(httpContext *gin.Context, body any) error {
	if err := httpContext.ShouldBindJSON(body); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func HandleJSONResponse(httpContext *gin.Context, response any, err error) {
	if err != nil {
		code, message := parseError(err)
		httpContext.AbortWithStatusJSON(http.StatusBadRequest, ApiError{
			Error:   true,
			Code:    code,
			Message: message,
		})
		return
	}

	httpContext.AbortWithStatusJSON(http.StatusOK, response)
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
