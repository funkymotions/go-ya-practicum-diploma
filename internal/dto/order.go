package dto

import "github.com/funkymotions/go-ya-practicum-diploma/internal/model"

type RegisterOrderRequest struct {
	OrderID string `json:"order_id"`
	UserID  uint   `json:"user_id"`
}

type GetUserOrderResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    *float64          `json:"accrual,omitempty"`
	UploadedAt string            `json:"uploaded_at"`
}
