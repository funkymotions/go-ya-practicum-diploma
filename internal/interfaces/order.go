package interfaces

import (
	"context"
	"sync"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type OrderService interface {
	RegisterOrder(*dto.RegisterOrderRequest) error
	GetUserOrders(userID uint) (*[]model.Order, error)
	EnqueueOrdersFromAccrualSystem(context.Context, *sync.WaitGroup) error
	ProcessOrderFromAccrual(order *model.Order) error
}

type OrderRepository interface {
	FindOneByID(id string) (*model.Order, error)
	CreateOrder(orderID string, userID uint) (*model.Order, error)
	GetUserOrders(userID uint) (*[]model.Order, error)
	GetUserOrderBalance(userID uint) (float64, error)
	FindManyByStatus(statuses ...model.OrderStatus) (*[]model.Order, error)
	UpdateOrder(order *model.Order) (*model.Order, error)
}
