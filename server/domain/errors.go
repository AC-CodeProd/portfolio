package domain

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeRequiredField ErrorCode = "REQUIRED_FIELD"
	ErrCodeInvalidFormat ErrorCode = "INVALID_FORMAT"
	ErrCodeInvalidLength ErrorCode = "INVALID_LENGTH"

	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"

	ErrCodeDatabase     ErrorCode = "DATABASE_ERROR"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeTimeout      ErrorCode = "TIMEOUT_ERROR"
	ErrCodeRateLimit    ErrorCode = "RATE_LIMIT_ERROR"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"
)

type DomainError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Field   string                 `json:"field,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	Cause   error                  `json:"cause,omitempty"`
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func (e *DomainError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidation, ErrCodeRequiredField, ErrCodeInvalidFormat, ErrCodeInvalidLength:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeUnauthorized, ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeDatabase, ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (e *DomainError) Title() string {
	switch e.Code {
	case ErrCodeValidation:
		return "Validation Error"
	case ErrCodeRequiredField:
		return "Required Field Missing"
	case ErrCodeInvalidFormat:
		return "Invalid Format"
	case ErrCodeInvalidLength:
		return "Invalid Length"
	case ErrCodeNotFound:
		return "Resource Not Found"
	case ErrCodeAlreadyExists:
		return "Resource Already Exists"
	case ErrCodeUnauthorized:
		return "Unauthorized"
	case ErrCodeForbidden:
		return "Forbidden"
	case ErrCodeDatabase:
		return "Database Error"
	case ErrCodeInternal:
		return "Internal Server Error"
	case ErrCodeTimeout:
		return "Request Timeout"
	case ErrCodeRateLimit:
		return "Rate Limit Exceeded"
	case ErrCodeTokenExpired:
		return "Token Expired"
	default:
		return "Unknown Error"
	}
}

func NewValidationError(message string, field string, cause *error) *DomainError {
	domainErr := &DomainError{
		Code:    ErrCodeValidation,
		Message: message,
		Field:   field,
	}

	if cause != nil {
		domainErr.Cause = *cause
	}

	return domainErr
}

func NewRequiredFieldError(field string) *DomainError {
	return &DomainError{
		Code:    ErrCodeRequiredField,
		Message: fmt.Sprintf("Field '%s' is required", field),
		Field:   field,
	}
}

func NewInvalidFormatError(field string, expectedFormat string) *DomainError {
	return &DomainError{
		Code:    ErrCodeInvalidFormat,
		Message: fmt.Sprintf("Field '%s' has invalid format, expected: %s", field, expectedFormat),
		Field:   field,
		Details: map[string]interface{}{
			"expected_format": expectedFormat,
		},
	}
}

func NewNotFoundError(resource string, identifier string) *DomainError {
	message := fmt.Sprintf("%s not found", resource)
	if identifier != "" {
		message = fmt.Sprintf("%s with identifier '%s' not found", resource, identifier)
	}

	return &DomainError{
		Code:    ErrCodeNotFound,
		Message: message,
		Details: map[string]interface{}{
			"resource":   resource,
			"identifier": identifier,
		},
	}
}

func NewAlreadyExistsError(resource string, identifier string) *DomainError {
	return &DomainError{
		Code:    ErrCodeAlreadyExists,
		Message: fmt.Sprintf("%s with identifier '%s' already exists", resource, identifier),
		Details: map[string]interface{}{
			"resource":   resource,
			"identifier": identifier,
		},
	}
}

func NewDatabaseError(operation string, cause error) *DomainError {
	return &DomainError{
		Code:    ErrCodeDatabase,
		Message: fmt.Sprintf("Database error during %s operation", operation),
		Cause:   cause,
		Details: map[string]interface{}{
			"operation": operation,
		},
	}
}

func NewInternalError(message string, cause error) *DomainError {
	return &DomainError{
		Code:    ErrCodeInternal,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeForbidden,
		Message: message,
	}
}

func NewRateLimitError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeRateLimit,
		Message: message,
	}
}

func NewTokenExpiredError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeTokenExpired,
		Message: message,
	}
}

func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

func AsDomainError(err error) (*DomainError, bool) {
	domainErr, ok := err.(*DomainError)
	return domainErr, ok
}
