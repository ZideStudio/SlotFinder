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
	ERR_USERNAME_MISSING       = err("USERNAME_MISSING", 0)
	// Event
	ERR_EVENT_NOT_FOUND          = err("EVENT_NOT_FOUND", http.StatusNotFound)
	ERR_EVENT_START_AFTER_END    = err("EVENT_START_AFTER_END", 0)
	ERR_EVENT_START_BEFORE_TODAY = err("EVENT_START_NOT_TODAY", 0)
	ERR_EVENT_DURATION_TOO_SHORT = err("EVENT_DURATION_TOO_SHORT", 0)
	ERR_EVENT_ALREADY_JOINED     = err("EVENT_ALREADY_JOINED", 0)
)

var CUSTOM_ERRORS = []CustomError{
	// Generic
	ERR_SERVER_ERROR,
	// Auth
	ERR_NOT_AUTHENTICATED,
	ERR_TOKEN_INVALID,
	ERR_TOKEN_EXPIRED,
	// Account
	ERR_INVALID_EMAIL_FORMAT,
	ERR_USERNAME_ALREADY_TAKEN,
	ERR_EMAIL_ALREADY_EXISTS,
	ERR_USERNAME_MISSING,
	// Event
	ERR_EVENT_NOT_FOUND,
	ERR_EVENT_START_AFTER_END,
	ERR_EVENT_START_BEFORE_TODAY,
	ERR_EVENT_DURATION_TOO_SHORT,
	ERR_EVENT_ALREADY_JOINED,
}

var CUSTOM_ERRORS_MAP = func() map[string]CustomError {
	m := make(map[string]CustomError)
	for _, ce := range CUSTOM_ERRORS {
		m[ce.Err.Error()] = ce
	}
	return m
}()
