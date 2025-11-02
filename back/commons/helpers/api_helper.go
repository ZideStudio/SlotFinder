package helpers

import (
	"app/commons/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code string `json:"code"`
}

func SetHttpContextBody(httpContext *gin.Context, body any) error {
	if err := httpContext.ShouldBindJSON(body); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func HandleJSONResponse(httpContext *gin.Context, response any, err error) {
	if err == nil {
		httpContext.AbortWithStatusJSON(http.StatusOK, response)
		return
	}

	code, status := parseError(err)
	httpContext.AbortWithStatusJSON(status, ApiError{
		Code: code,
	})
}

func parseError(err error) (code string, status int) {
	if customError, exists := constants.CUSTOM_ERRORS_MAP[err.Error()]; exists {
		code = err.Error()

		if customError.StatusCode != nil {
			status = *customError.StatusCode
		} else {
			status = http.StatusBadRequest
		}

		return
	}

	code = constants.ERR_SERVER_ERROR.Err.Error()
	status = http.StatusInternalServerError

	return
}
