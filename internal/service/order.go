package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
)

type accuralServiceOrderResponse struct {
	Order   string            `json:"order"`
	Status  model.OrderStatus `json:"status"`
	Accrual float64           `json:"accrual"`
}

type orderService struct {
	repo          interfaces.OrderRepository
	accrualClient interfaces.RESTClient
	queue         chan *model.Order
	semaphore     interfaces.Semaphore
	logger        interfaces.Logger
}

type OrderServiceConf struct {
	OrderRepo      interfaces.OrderRepository
	AccrualClient  interfaces.RESTClient
	OrderQueue     chan *model.Order
	RequestLimiter interfaces.Semaphore
	Logger         interfaces.Logger
}

func NewOrderService(
	conf *OrderServiceConf,
) *orderService {
	return &orderService{
		repo:          conf.OrderRepo,
		accrualClient: conf.AccrualClient,
		queue:         conf.OrderQueue,
		semaphore:     conf.RequestLimiter,
		logger:        conf.Logger,
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

// TODO: implement semaphore or limit number of goroutines
func (s *orderService) EnqueueOrdersFromAccrualSystem(ctx context.Context, wg *sync.WaitGroup) error {
	orders, err := s.repo.FindManyByStatus(model.StatusNew, model.StatusProcessing)
	if err != nil {
		s.logger.Errorf("error fetching orders from repository: %v", err)
		return err
	}
	s.logger.Infof("found orders to process: %d", len(*orders))
	// Call accrual service for each order and fill in the workers queue
	for _, order := range *orders {
		wg.Add(1)
		go s.callAccrualService(ctx, wg, order)
	}
	return nil
}

func (s *orderService) callAccrualService(ctx context.Context, wg *sync.WaitGroup, order model.Order) {
	defer wg.Done()
	// limit number of concurrent requests to remote service using semaphore
	s.semaphore.Acquire()
	defer s.semaphore.Release()
	s.logger.Infof("requesting order with ID = %s", order.ID)
	select {
	case <-ctx.Done():
		return
	default:
		u, err := url.Parse("api/orders/")
		if err != nil {
			s.logger.Errorf("error parsing URL: %v", err)
			return
		}
		resp, err := s.accrualClient.Get(u.JoinPath(order.ID).Path)
		if err != nil {
			s.logger.Errorf("error calling accrual service: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			s.logger.Warnf("accrual service returned non-200 status: %d", resp.StatusCode)
			return
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Errorf("error reading response body: %v", err)
			return
		}
		var orderResp accuralServiceOrderResponse
		err = json.Unmarshal(data, &orderResp)
		if err != nil {
			s.logger.Errorf("error unmarshalling response body: %v", err)
			return
		}
		if orderResp.Status == model.StatusInvalid {
			order.Accrual = nil
		} else {
			order.Accrual = &orderResp.Accrual
		}
		order.OrderStatus = orderResp.Status
		s.queue <- &order
		// return
	}
}

func (s *orderService) ProcessOrderFromAccrual(order *model.Order) error {
	if order.OrderStatus == model.StatusInvalid || order.OrderStatus == model.StatusProcessed {
		_, err := s.repo.UpdateOrder(order)
		if err != nil {
			return err
		}
	}
	return nil
}
