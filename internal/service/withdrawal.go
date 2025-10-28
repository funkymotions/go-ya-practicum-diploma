package service

import (
	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type withdrawalService struct {
	repo interfaces.WithdrawalRepository
}

func NewWithdrawalService(r interfaces.WithdrawalRepository) *withdrawalService {
	return &withdrawalService{
		repo: r,
	}
}

func (s *withdrawalService) GetUserWithdrawals(userID uint) (*[]model.Withdrawal, error) {
	withdrawals, err := s.repo.GetUserWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	if len(*withdrawals) == 0 {
		return nil, apperrors.ErrNoWithdrawals
	}
	return withdrawals, nil
}
