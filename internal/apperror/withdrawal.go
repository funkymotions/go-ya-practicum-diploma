package apperrors

import "net/http"

var ErrNoWithdrawals = &AppError{
	Message:    "no withdrawals found for user",
	StatusCode: http.StatusNoContent,
}
