package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorCode string

const (
	ErrCode_BadRequest          ErrorCode = "BAD_REQUEST"
	ErrCode_Forbidden           ErrorCode = "FORBIDDEN"
	ErrCode_ResourceNotFound    ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCode_InternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCode_Unregistered        ErrorCode = "UNREGISTERED_ERRCODE"
)

var (
	ErrorMappings = map[ErrorCode]ErrorMapping{
		ErrCode_BadRequest:          {HTTPCode: http.StatusBadRequest, Message: "request has invalid parameter(s) or header(s)."},
		ErrCode_Forbidden:           {HTTPCode: http.StatusForbidden, Message: "the operation is forbidden."},
		ErrCode_ResourceNotFound:    {HTTPCode: http.StatusNotFound, Message: "resource not found."},
		ErrCode_InternalServerError: {HTTPCode: http.StatusInternalServerError, Message: "There is a problem on our end. Please try again later."},
	}
)

type ErrorMapping struct {
	HTTPCode int
	Message  string
}

type ServiceError struct {
	HTTPStatus int           `json:"-"`
	Code       string        `json:"code,omitempty"`
	Message    string        `json:"message,omitempty"`
	Errors     []ErrorDetail `json:"error,omitempty"`
}

// Error implements the error interface
func (e ServiceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

type ErrorDetail struct {
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
	Path      string `json:"path,omitempty"`
}

func NewServiceError(errCode ErrorCode, errDetails []ErrorDetail) ServiceError {
	mapping, ok := ErrorMappings[errCode]
	if !ok {
		mapping = ErrorMappings[ErrCode_Unregistered]
	}

	return ServiceError{
		HTTPStatus: mapping.HTTPCode,
		Code:       string(errCode),
		Message:    mapping.Message,
		Errors:     errDetails,
	}
}

// ParseValidationErrors converts validator.ValidationErrors to []ErrorDetail
func ParseValidationErrors(err error) []ErrorDetail {
	var errorDetails []ErrorDetail
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, fieldErr := range validationErrs {
			errorDetails = append(errorDetails, ErrorDetail{
				ErrorCode: "VALIDATION_ERROR",
				Message:   fmt.Sprintf("Field '%s' failed validation: %s", fieldErr.Field(), fieldErr.Tag()),
				Path:      fieldErr.Namespace(),
			})
		}
	}

	return errorDetails
}
