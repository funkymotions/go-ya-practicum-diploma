package apperrors

import "net/http"

var ErrAccountInsufficientFunds = &AppError{
	StatusCode: http.StatusPaymentRequired,
	Message:    "Insufficient funds in account",
}
