package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ApiError struct {
	Error   bool   `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func extractFieldNameFromError(errString string, obj any) (string, *string) {
	if strings.Contains(errString, "is not set") {
		errParts := strings.Split(errString, " ")
		fieldName := strings.Trim(errParts[1], "'")

		s := reflect.ValueOf(obj).Elem()
		for i := 0; i < s.NumField(); i++ {
			typeField := s.Type().Field(i)
			jsonTag := typeField.Tag.Get("json")
			if jsonTag == fieldName {
				return fieldName, &fieldName
			}
		}
	}

	return "message", nil
}

func extractFieldName(fullFieldName string) string {
	parts := strings.Split(fullFieldName, ".")
	fieldName := parts[len(parts)-1]
	return strings.ToLower(string(rune(fieldName[0]))) + fieldName[1:]
}

func parseTypeName(exceptedTypeName string) string {
	if strings.Contains(exceptedTypeName, ".") {
		typeName := strings.Split(exceptedTypeName, ".")[1]
		return typeName
	}

	return exceptedTypeName
}

func ShouldBindJSON(c *gin.Context, obj any) (err error) {
	if err = c.ShouldBindJSON(&obj); err != nil {
		errMap := make(map[string]string)
		var field *string
		errContentLower := strings.ToLower(err.Error())

		if errParam, ok := err.(*json.UnmarshalTypeError); ok {
			fieldName := errParam.Field
			expectedType := errParam.Type.String()
			fieldExtracted := extractFieldName(fieldName)
			field = &fieldExtracted
			errMap[fieldExtracted] = "must be of the type " + parseTypeName(expectedType)
		} else if strings.Contains(errContentLower, "is not set") {
			fieldExtracted, fieldNameValue := extractFieldNameFromError(err.Error(), obj)
			field = fieldNameValue
			errMap[fieldExtracted] = "can't be empty"
		} else if strings.Contains(errContentLower, "field validation for") && strings.Contains(errContentLower, "failed on the 'required' tag") {
			fieldName := strings.Split(err.Error(), "'")[1]
			fieldExtracted := extractFieldName(fieldName)
			field = &fieldExtracted
			errMap[fieldExtracted] = "can't be empty"
		} else if strings.Contains(errContentLower, "parsing time") {
			fieldExtracted, fieldNameValue := extractFieldNameFromError(err.Error(), obj)
			field = fieldNameValue
			errMap[fieldExtracted] = "one of the time format is invalid, please use the 2006-01-02T15:04:05Z07:00 format"
		} else {
			errMap["message"] = "one of the field is invalid"
		}

		var response ApiError
		if field != nil {
			response = ApiError{
				Error:   true,
				Code:    "INVALID_INPUT",
				Message: fmt.Sprintf("%s: %s", *field, errMap[*field]),
			}
		} else {
			response = ApiError{
				Error:   true,
				Code:    "INVALID_INPUT",
				Message: errMap["message"],
			}
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	return
}

func ShouldBindForm(c *gin.Context, obj any) (err error) {
	if err = c.ShouldBindWith(obj, binding.FormMultipart); err != nil {
		code, message := parseError(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiError{
			Error:   true,
			Code:    code,
			Message: message,
		})
	}

	return
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
