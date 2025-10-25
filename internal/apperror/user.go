package apperrors

import (
	"net/http"
)

var ErrUserAlreadyExists = &AppError{
	Message:    "user already exists",
	StatusCode: http.StatusConflict,
}

var ErrUserNotFound = &AppError{
	Message:    "user not found",
	StatusCode: http.StatusNotFound,
}

var ErrInvalidCredentials = &AppError{
	Message:    "invalid credentials",
	StatusCode: http.StatusUnauthorized,
}
