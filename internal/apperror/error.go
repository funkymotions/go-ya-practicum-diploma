package apperrors

import "fmt"

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("App error: %s (status code: %d)", e.Message, e.StatusCode)
}
