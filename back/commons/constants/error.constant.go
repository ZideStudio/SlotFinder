package constants

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Error      error // Custom error code
	StatusCode *int  // Http Status Code, keep nil for default status code 400
}

func err(code string, httpStatus int) CustomError {
	var httpStatusPtr *int
	if httpStatus != 0 {
		httpStatusPtr = &httpStatus
	}

	return CustomError{errors.New(code), httpStatusPtr}
}

var (
	// Generic
	ERR_SERVER_ERROR = err("SERVER_ERROR", 0)
	// Auth
	ERR_NO_COOKIE     = err("NO_COOKIE", http.StatusUnauthorized)
	ERR_TOKEN_INVALID = err("TOKEN_INVALID", http.StatusUnauthorized)
	ERR_TOKEN_EXPIRED = err("TOKEN_EXPIRED", http.StatusUnauthorized)
	// Account
	ERR_INVALID_EMAIL_FORMAT   = err("INVALID_EMAIL_FORMAT", 0)
	ERR_USERNAME_ALREADY_TAKEN = err("USERNAME_ALREADY_TAKEN", 0)
	ERR_EMAIL_ALREADY_EXISTS   = err("EMAIL_ALREADY_EXISTS", 0)
	// Event
	ERR_EVENT_START_AFTER_END = err("EVENT_START_AFTER_END", 0)
	ERR_EVENT_START_NOT_TODAY = err("EVENT_START_NOT_TODAY", 0)
)

var CUSTOM_ERRORS = []error{
	// Generic
	ERR_SERVER_ERROR.Error,
	// Auth
	ERR_NO_COOKIE.Error,
	ERR_TOKEN_INVALID.Error,
	ERR_TOKEN_EXPIRED.Error,
	// Account
	ERR_INVALID_EMAIL_FORMAT.Error,
	ERR_USERNAME_ALREADY_TAKEN.Error,
	ERR_EMAIL_ALREADY_EXISTS.Error,
	// Event
	ERR_EVENT_START_AFTER_END.Error,
	ERR_EVENT_START_NOT_TODAY.Error,
}
