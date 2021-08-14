package errorhandler

import (
	"fmt"
	"net/http"
)

type ErrorCode string

type (
	ServiceError interface {
		error
		ErrorCode() string
		ToMap() map[string]interface{}
	}

	Error struct {
		error
		Code        ErrorCode
		Value       string
		Description string
	}
)

func (e Error) Error() string {
	if e.error != nil {
		return fmt.Sprintf("[%s] %s", e.ErrorCode(), e.error.Error())
	}
	if e.Description != "" {
		return fmt.Sprintf("[%s] %s", e.ErrorCode(), e.Description)
	}
	return fmt.Sprintf("[%s] %s", e.ErrorCode(), e.Value)
}

func (e Error) ErrorCode() string {
	return string(e.Code)
}

// ToMap converts Error object to map[string]interface{}
func (e Error) ToMap() map[string]interface{} {
	r := map[string]interface{}{
		"code": string(e.Code),
	}

	if string(e.Value) != "" {
		r["value"] = string(e.Value)
	}

	if string(e.Description) != "" {
		r["description"] = string(e.Description)
	}

	return r
}

// GetHttpCode return a Http error code
func (e Error) GetHttpCode() int {
	switch e.Code {
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrService:
		return http.StatusInternalServerError
	case ErrNotAllowed:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusBadRequest
	}
}

// New creates an Error object
func NewError(code ErrorCode, value ...string) *Error {
	e := &Error{Code: code}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// NewWithDesc creates an Error object with description
func NewErrorWithDesc(code ErrorCode, desc string, value ...string) *Error {
	e := &Error{Code: code, Description: desc}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// FromError creates a new Error (ErrService) from common golang error
func FromVanillaError(err error) *Error {
	return &Error{
		Code:  ErrService,
		Value: "",
	}
}
