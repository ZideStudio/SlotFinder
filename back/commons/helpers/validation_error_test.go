package helpers

import (
	"app/commons/constants"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validationSampleBody struct {
	Name     string    `json:"name" binding:"required,min=3,max=5"`
	Email    string    `json:"email" binding:"omitempty,email"`
	Count    *int      `json:"count" binding:"omitempty,min=1,max=10"`
	Language string    `json:"language" binding:"omitempty,oneof=en fr"`
	Accepted *bool     `json:"accepted" binding:"omitempty,eq=true"`
	StartsAt time.Time `json:"startsAt"`
}

func bindSampleBody(t *testing.T, body string) *ValidationError {
	t.Helper()
	c, _ := newHelperTestContext([]byte(body))

	var data validationSampleBody
	err := c.ShouldBindJSON(&data)
	require.Error(t, err)

	return ParseBindingError(err)
}

func assertSingleFieldError(t *testing.T, validationError *ValidationError, expected FieldError) {
	t.Helper()
	assert.Equal(t, constants.ERR_INVALID_BODY, validationError.Code)
	assert.Equal(t, []FieldError{expected}, validationError.Fields)
}

func TestParseBindingError_Required(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "name", Code: "FIELD_REQUIRED"})
}

func TestParseBindingError_StringTooShort(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"ab"}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "name", Code: "FIELD_TOO_SHORT"})
}

func TestParseBindingError_StringTooLong(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abcdef"}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "name", Code: "FIELD_TOO_LONG"})
}

func TestParseBindingError_NumberTooSmall(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","count":0}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "count", Code: "FIELD_TOO_SMALL"})
}

func TestParseBindingError_NumberTooLarge(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","count":11}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "count", Code: "FIELD_TOO_LARGE"})
}

func TestParseBindingError_InvalidEmail(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","email":"not-an-email"}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "email", Code: "FIELD_INVALID_EMAIL"})
}

func TestParseBindingError_OneOf(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","language":"de"}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "language", Code: "FIELD_INVALID_VALUE"})
}

func TestParseBindingError_Eq(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","accepted":false}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "accepted", Code: "FIELD_INVALID_VALUE"})
}

func TestParseBindingError_MultipleFields(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"ab","email":"nope"}`)
	assert.Equal(t, constants.ERR_INVALID_BODY, validationError.Code)
	assert.Equal(t, []FieldError{
		{Field: "name", Code: "FIELD_TOO_SHORT"},
		{Field: "email", Code: "FIELD_INVALID_EMAIL"},
	}, validationError.Fields)
}

func TestParseBindingError_InvalidType(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","count":"three"}`)
	assertSingleFieldError(t, validationError, FieldError{Field: "count", Code: "FIELD_INVALID_TYPE"})
}

func TestParseBindingError_InvalidDate(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":"abc","startsAt":"not-a-date"}`)
	assertSingleFieldError(t, validationError, FieldError{Code: "FIELD_INVALID_FORMAT"})
}

func TestParseBindingError_MalformedJSON(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `{"name":`)
	assert.Equal(t, constants.ERR_INVALID_JSON, validationError.Code)
	assert.Empty(t, validationError.Fields)
}

func TestParseBindingError_SyntaxError(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, `not-json`)
	assert.Equal(t, constants.ERR_INVALID_JSON, validationError.Code)
	assert.Empty(t, validationError.Fields)
}

func TestParseBindingError_EmptyBody(t *testing.T) {
	t.Parallel()
	validationError := bindSampleBody(t, ``)
	assert.Equal(t, constants.ERR_INVALID_JSON, validationError.Code)
	assert.Empty(t, validationError.Fields)
}

func TestParseBindingError_UnknownError(t *testing.T) {
	t.Parallel()
	validationError := ParseBindingError(errors.New("unexpected"))
	assert.Equal(t, constants.ERR_INVALID_BODY, validationError.Code)
	assert.Empty(t, validationError.Fields)
}

func TestNewFieldError_SingleField(t *testing.T) {
	t.Parallel()
	validationError := NewFieldError(constants.ERR_INVALID_PASSWORD_FORMAT, "password")
	assertSingleFieldError(t, validationError, FieldError{Field: "password", Code: "INVALID_PASSWORD_FORMAT"})
}

func TestNewFieldError_MultipleFields(t *testing.T) {
	t.Parallel()
	validationError := NewFieldError(constants.ERR_AVAILABILITY_INVALID_TIME_INTERVAL, "startsAt", "endsAt")
	assert.Equal(t, constants.ERR_INVALID_BODY, validationError.Code)
	assert.Equal(t, []FieldError{
		{Field: "startsAt", Code: "AVAILABILITY_INVALID_TIME_INTERVAL"},
		{Field: "endsAt", Code: "AVAILABILITY_INVALID_TIME_INTERVAL"},
	}, validationError.Fields)
}

func TestFieldErrors_EmptyReturnsNil(t *testing.T) {
	t.Parallel()
	var fieldErrors FieldErrors
	err := fieldErrors.Err()
	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestFieldErrors_KeepsOrderAndDuplicateFields(t *testing.T) {
	t.Parallel()
	var fieldErrors FieldErrors
	fieldErrors.Add(constants.ERR_EVENT_START_AFTER_END, "endsAt")
	fieldErrors.Add(constants.ERR_EVENT_START_BEFORE_TODAY, "startsAt")
	fieldErrors.Add(constants.ERR_EVENT_DURATION_TOO_SHORT, "endsAt")

	assert.Equal(t, &ValidationError{Code: constants.ERR_INVALID_BODY, Fields: []FieldError{
		{Field: "endsAt", Code: "EVENT_START_AFTER_END"},
		{Field: "startsAt", Code: "EVENT_START_NOT_TODAY"},
		{Field: "endsAt", Code: "EVENT_DURATION_TOO_SHORT"},
	}}, fieldErrors.Err())
}

func TestFieldErrors_HasCollectedField(t *testing.T) {
	t.Parallel()
	var fieldErrors FieldErrors
	fieldErrors.Add(constants.ERR_INVALID_EMAIL_FORMAT, "email")
	assert.True(t, fieldErrors.Has("email"))
}

func TestFieldErrors_HasNotCollectedField(t *testing.T) {
	t.Parallel()
	var fieldErrors FieldErrors
	fieldErrors.Add(constants.ERR_INVALID_EMAIL_FORMAT, "email")
	assert.False(t, fieldErrors.Has("password"))
}
