package apierror

import "fmt"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(statusCode int, code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func Wrap(err error, statusCode int, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    err.Error(),
		StatusCode: statusCode,
	}
}

var (
	ErrBadRequest         = New(400, "BAD_REQUEST", "Bad request")
	ErrUnauthorized       = New(401, "UNAUTHORIZED", "Unauthorized")
	ErrForbidden          = New(403, "FORBIDDEN", "Forbidden")
	ErrNotFound           = New(404, "NOT_FOUND", "Resource not found")
	ErrInternalServer     = New(500, "INTERNAL_SERVER_ERROR", "Internal server error")
	ErrInvalidCredentials = New(401, "INVALID_CREDENTIALS", "Invalid email or password")
	ErrEmailAlreadyExists = New(400, "EMAIL_EXISTS", "Email already registered")
	ErrInvalidLanguage    = New(400, "INVALID_LANGUAGE", "Invalid or unsupported language")
)
