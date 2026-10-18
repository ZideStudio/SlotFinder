package helpers

import (
	"app/commons/constants"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ApiError struct {
	Code   string       `json:"code"`
	Fields []FieldError `json:"fields,omitempty"`
}

func SetHttpContextBody(httpContext *gin.Context, body any) error {
	if err := httpContext.ShouldBindJSON(body); err != nil {
		validationError := parseBodyBindingError(err, body)
		HandleJSONResponse(httpContext, nil, validationError)
		return validationError
	}

	return nil
}

func HandleJSONResponse(httpContext *gin.Context, response any, err error) {
	if err == nil {
		httpContext.AbortWithStatusJSON(http.StatusOK, response)
		return
	}

	var validationError *ValidationError
	if errors.As(err, &validationError) {
		httpContext.AbortWithStatusJSON(http.StatusBadRequest, ApiError{
			Code:   validationError.Code.Err.Error(),
			Fields: validationError.Fields,
		})
		return
	}

	code, status := parseError(err)

	if code == constants.ERR_SERVER_ERROR.Err.Error() {
		log.Error().Err(err).Msg("Internal server error")
	}

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
