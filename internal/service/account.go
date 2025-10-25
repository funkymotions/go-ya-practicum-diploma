package service

import (
	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
)

type accountService struct {
	orderRepo      orderRepository
	withdrawalRepo withdrawalRepository
}

func NewAccountService(o orderRepository, w withdrawalRepository) *accountService {
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
	return balance - withdrawn, withdrawn, nil
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
