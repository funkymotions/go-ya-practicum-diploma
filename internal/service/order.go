package service

import (
	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
)

type orderRepository interface {
	FindOneByID(id string) (*model.Order, error)
	CreateOrder(orderID string, userID uint) (*model.Order, error)
	GetUserOrders(userID uint) (*[]model.Order, error)
	GetUserOrderBalance(userID uint) (float64, error)
}

type orderService struct {
	repo orderRepository
}

func NewOrderService(r orderRepository) *orderService {
	return &orderService{
		repo: r,
	}
}

func (s *orderService) RegisterOrder(input *dto.RegisterOrderRequest) error {
	if !utils.IsValidLuhn(input.OrderID) {
		return apperrors.ErrOrderInvalidID
	}
	order, err := s.repo.FindOneByID(input.OrderID)
	if err != nil && err != apperrors.ErrOrderNotFound {
		return err
	}
	if order != nil && order.UserID != input.UserID {
		return apperrors.ErrOrderAssignedWithUser
	}
	_, err = s.repo.CreateOrder(input.OrderID, input.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *orderService) GetUserOrders(userID uint) (*[]model.Order, error) {
	orders, err := s.repo.GetUserOrders(userID)
	if err != nil {
		return nil, err
	}
	if len(*orders) == 0 {
		return nil, apperrors.ErrOrderListEmpty
	}
	return orders, nil
}
