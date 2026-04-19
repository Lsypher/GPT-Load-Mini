package errors

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNoActiveKeys   = errors.New("no active keys available")
	ErrInternalServer = errors.New("internal server error")
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code int, msg string) *APIError {
	return &APIError{Code: code, Message: msg}
}
