package apperror

import "fmt"

type Code string

const (
	CodeNotFound      Code = "not_found"
	CodeConflict      Code = "conflict"
	CodeValidation    Code = "validation"
	CodeUnauthorized  Code = "unauthorized"
	CodeForbidden     Code = "forbidden"
	CodeInternal      Code = "internal"
	CodeInvalidState  Code = "invalid_state"
	CodeUnavailable   Code = "unavailable"
)

type AppError struct {
	Code    Code
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

func New(code Code, msg string, cause error) *AppError {
	return &AppError{Code: code, Message: msg, Cause: cause}
}

