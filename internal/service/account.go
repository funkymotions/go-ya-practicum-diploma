package service

import (
	"math"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
)

type accountService struct {
	orderRepo      interfaces.OrderRepository
	withdrawalRepo interfaces.WithdrawalRepository
}

func NewAccountService(o interfaces.OrderRepository, w interfaces.WithdrawalRepository) *accountService {
	return &accountService{
		orderRepo:      o,
		withdrawalRepo: w,
	}
}

func (s *accountService) CalculateAccountBalance(userID uint) (float64, float64, error) {
	balance, err := s.orderRepo.GetUserOrderBalance(userID)
	if err != nil {
		return 0, 0, err
	}
	withdrawn, err := s.withdrawalRepo.GetUserWithdrawalBalance(userID)
	if err != nil {
		return 0, 0, err
	}
	roundWithdrawn := math.Round(withdrawn*1000) / 1000
	roundBalance := math.Round((balance-roundWithdrawn)*1000) / 1000
	return roundBalance, roundWithdrawn, nil
}

func (s *accountService) WithdrawFromAccount(userID uint, data *dto.AccountBalanceWithdrawRequest) error {
	if !utils.IsValidLuhn(data.Order) {
		return apperrors.ErrOrderInvalidID
	}
	balance, _, err := s.CalculateAccountBalance(userID)
	if err != nil {
		return err
	}
	if balance < data.Sum {
		return apperrors.ErrAccountInsufficientFunds
	}
	err = s.withdrawalRepo.CreateWithdrawal(userID, data.Order, data.Sum)
	if err != nil {
		return err
	}
	return nil
}
