package helpers

import (
	"app/commons/constants"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code"`
	Param string `json:"param,omitempty"`
}

type ValidationError struct {
	Code   constants.CustomError
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return e.Code.Err.Error()
}

func init() {
	// Report field names as they appear in the JSON payload
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "" || name == "-" {
				return field.Name
			}
			return name
		})
	}
}

// ParseBindingError converts an error returned by gin binding into a ValidationError
func ParseBindingError(err error) *ValidationError {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make([]FieldError, 0, len(validationErrors))
		for _, fieldError := range validationErrors {
			fields = append(fields, FieldError{
				Field: fieldError.Field(),
				Code:  validationTagToError(fieldError).Err.Error(),
				Param: fieldError.Param(),
			})
		}
		return &ValidationError{Code: constants.ERR_INVALID_BODY, Fields: fields}
	}

	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		return &ValidationError{
			Code:   constants.ERR_INVALID_BODY,
			Fields: []FieldError{{Field: typeError.Field, Code: constants.ERR_FIELD_INVALID_TYPE.Err.Error()}},
		}
	}

	var timeError *time.ParseError
	if errors.As(err, &timeError) {
		return &ValidationError{
			Code:   constants.ERR_INVALID_BODY,
			Fields: []FieldError{{Code: constants.ERR_FIELD_INVALID_FORMAT.Err.Error()}},
		}
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return &ValidationError{Code: constants.ERR_INVALID_JSON}
	}

	return &ValidationError{Code: constants.ERR_INVALID_BODY}
}

// parseBodyBindingError validates the fields decoded around a type error
func parseBodyBindingError(err error, body any) *ValidationError {
	validationError := ParseBindingError(err)

	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		return validationError
	}

	// every invalid field is reported at once
	validateErr := binding.Validator.ValidateStruct(body)
	if validateErr == nil {
		return validationError
	}

	for _, fieldError := range ParseBindingError(validateErr).Fields {
		if fieldError.Field != typeError.Field {
			validationError.Fields = append(validationError.Fields, fieldError)
		}
	}

	return validationError
}

func validationTagToError(fieldError validator.FieldError) constants.CustomError {
	switch fieldError.Tag() {
	case "required":
		return constants.ERR_FIELD_REQUIRED
	case "email":
		return constants.ERR_FIELD_INVALID_EMAIL
	case "min":
		if hasLength(fieldError.Kind()) {
			return constants.ERR_FIELD_TOO_SHORT
		}
		return constants.ERR_FIELD_TOO_SMALL
	case "max":
		if hasLength(fieldError.Kind()) {
			return constants.ERR_FIELD_TOO_LONG
		}
		return constants.ERR_FIELD_TOO_LARGE
	default:
		return constants.ERR_FIELD_INVALID_VALUE
	}
}

func hasLength(kind reflect.Kind) bool {
	return kind == reflect.String || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Array
}
