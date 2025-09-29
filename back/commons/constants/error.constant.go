package constants

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Err        error // Custom error code
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
	ERR_NOT_AUTHENTICATED          = err("NOT_AUTHENTICATED", http.StatusUnauthorized)
	ERR_TOKEN_INVALID              = err("TOKEN_INVALID", http.StatusUnauthorized)
	ERR_TOKEN_EXPIRED              = err("TOKEN_EXPIRED", http.StatusUnauthorized)
	ERR_PROVIDER_CONNECTION_FAILED = err("PROVIDER_CONNECTION_FAILED", 0)
	// Account
	ERR_INVALID_EMAIL_FORMAT   = err("INVALID_EMAIL_FORMAT", 0)
	ERR_USERNAME_ALREADY_TAKEN = err("USERNAME_ALREADY_TAKEN", 0)
	ERR_EMAIL_ALREADY_EXISTS   = err("EMAIL_ALREADY_EXISTS", 0)
	// Event
	ERR_EVENT_START_AFTER_END    = err("EVENT_START_AFTER_END", 0)
	ERR_EVENT_START_BEFORE_TODAY = err("EVENT_START_NOT_TODAY", 0)
)

var CUSTOM_ERRORS = []error{
	// Generic
	ERR_SERVER_ERROR.Err,
	// Auth
	ERR_NOT_AUTHENTICATED.Err,
	ERR_TOKEN_INVALID.Err,
	ERR_TOKEN_EXPIRED.Err,
	// Account
	ERR_INVALID_EMAIL_FORMAT.Err,
	ERR_USERNAME_ALREADY_TAKEN.Err,
	ERR_EMAIL_ALREADY_EXISTS.Err,
	// Event
	ERR_EVENT_START_AFTER_END.Err,
	ERR_EVENT_START_BEFORE_TODAY.Err,
}
