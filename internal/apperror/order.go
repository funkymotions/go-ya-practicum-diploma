package apperrors

import "net/http"

var ErrOrderAssignedWithUser = &AppError{
	Message:    "order already placed by another user",
	StatusCode: http.StatusConflict,
}

var ErrOrderAlreadyRegistered = &AppError{
	Message:    "order already registered",
	StatusCode: http.StatusOK,
}

var ErrOrderInvalidID = &AppError{
	Message:    "invalid order ID",
	StatusCode: http.StatusUnprocessableEntity,
}

var ErrOrderNotFound = &AppError{
	Message:    "order not found",
	StatusCode: http.StatusNotFound,
}

var ErrOrderListEmpty = &AppError{
	Message:    "no orders found for user",
	StatusCode: http.StatusNoContent,
}
