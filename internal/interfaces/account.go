package interfaces

import "github.com/funkymotions/go-ya-practicum-diploma/internal/dto"

type AccountService interface {
	CalculateAccountBalance(userID uint) (float64, float64, error)
	WithdrawFromAccount(userID uint, data *dto.AccountBalanceWithdrawRequest) error
}
