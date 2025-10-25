package service

import (
	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type withdrawalRepository interface {
	CreateWithdrawal(userID uint, orderID string, amount float64) error
	GetUserWithdrawalBalance(userID uint) (float64, error)
	GetUserWithdrawals(userID uint) (*[]model.Withdrawal, error)
}

type withdrawalService struct {
	repo withdrawalRepository
}

func NewWithdrawalService(r withdrawalRepository) *withdrawalService {
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
