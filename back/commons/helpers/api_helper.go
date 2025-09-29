package helpers

import (
	"app/commons/constants"
	"app/commons/lib"
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
	if lib.InArray(err, constants.CUSTOM_ERRORS) != -1 {
		code = err.Error()
		status = http.StatusBadRequest
		return
	}

	code = constants.ERR_SERVER_ERROR.Err.Error()
	status = http.StatusInternalServerError

	return
}
