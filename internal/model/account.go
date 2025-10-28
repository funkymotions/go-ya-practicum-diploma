package model

import "time"

type Account struct {
	ID        uint
	UserID    uint
	Balance   *float64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
