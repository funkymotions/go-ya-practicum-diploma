package model

import "time"

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          string
	UserID      uint
	OrderStatus OrderStatus
	Accrual     *float64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
