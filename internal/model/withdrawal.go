package model

import "time"

type Withdrawal struct {
	ID          uint
	UserID      uint
	OrderID     string
	Amount      float64
	ProcessedAt *time.Time
}
