package interfaces

import "github.com/funkymotions/go-ya-practicum-diploma/internal/model"

type WithdrawalRepository interface {
	CreateWithdrawal(userID uint, orderID string, amount float64) error
	GetUserWithdrawalBalance(userID uint) (float64, error)
	GetUserWithdrawals(userID uint) (*[]model.Withdrawal, error)
}

type WithdrawalService interface {
	GetUserWithdrawals(userID uint) (*[]model.Withdrawal, error)
}
