package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderTestSuite struct {
	suite.Suite
	repository *mocks.MockOrderRepository
	service    *orderService
	restClient *mocks.MockRESTClient
	semaphore  *mocks.MockSemaphore
	logger     *mocks.MockLogger
}

func (s *OrderTestSuite) SetupTest() {
	s.repository = mocks.NewMockOrderRepository(gomock.NewController(s.T()))
	s.restClient = mocks.NewMockRESTClient(gomock.NewController(s.T()))
	queueChan := make(chan *model.Order, 10)
	s.semaphore = mocks.NewMockSemaphore(gomock.NewController(s.T()))
	s.logger = mocks.NewMockLogger(gomock.NewController(s.T()))
	s.service = NewOrderService(
		&OrderServiceConf{
			OrderRepo:      s.repository,
			AccrualClient:  s.restClient,
			OrderQueue:     queueChan,
			RequestLimiter: s.semaphore,
			Logger:         s.logger,
		})
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) TestRegisterOrder() {
	type testCase struct {
		name         string
		mockBehavior func()
		wantErr      bool
		serviceArg   *dto.RegisterOrderRequest
	}
	testCases := []testCase{
		{
			name: "OK",
			serviceArg: &dto.RegisterOrderRequest{
				OrderID: "79927398713",
				UserID:  1,
			},
			mockBehavior: func() {
				s.repository.EXPECT().
					FindOneByID(gomock.Any()).
					Return(nil, apperrors.ErrOrderNotFound)
				s.repository.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid orderID",
			serviceArg: &dto.RegisterOrderRequest{
				OrderID: "799273987130",
				UserID:  1,
			},
			mockBehavior: nil,
			wantErr:      true,
		},
		{
			name: "repository unexpected error",
			serviceArg: &dto.RegisterOrderRequest{
				OrderID: "79927398713",
				UserID:  1,
			},
			mockBehavior: func() {
				s.repository.EXPECT().
					FindOneByID(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "Another user order",
			serviceArg: &dto.RegisterOrderRequest{
				OrderID: "79927398713",
				UserID:  2,
			},
			mockBehavior: func() {
				s.repository.EXPECT().
					FindOneByID(gomock.Any()).
					Return(&model.Order{
						ID:     "79927398713",
						UserID: 1,
					}, nil)
			},
			wantErr: true,
		},
		{
			name: "repository CreateOrder error",
			serviceArg: &dto.RegisterOrderRequest{
				OrderID: "79927398713",
				UserID:  1,
			},
			mockBehavior: func() {
				s.repository.EXPECT().
					FindOneByID(gomock.Any()).
					Return(nil, apperrors.ErrOrderNotFound)
				s.repository.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			err := s.service.RegisterOrder(tc.serviceArg)
			if tc.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *OrderTestSuite) TestGetUserOrders() {
	type testCase struct {
		name         string
		userID       uint
		mockBehavior func()
		wantErr      bool
	}
	testCases := []testCase{
		{
			name:   "OK",
			userID: 1,
			mockBehavior: func() {
				s.repository.EXPECT().
					GetUserOrders(gomock.Any()).
					Return(&[]model.Order{
						{ID: "79927398713", UserID: 1},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "No orders",
			userID: 1,
			mockBehavior: func() {
				s.repository.EXPECT().
					GetUserOrders(gomock.Any()).
					Return(&[]model.Order{}, nil)
			},
			wantErr: true,
		},
		{
			name:   "repository error",
			userID: 1,
			mockBehavior: func() {
				s.repository.EXPECT().
					GetUserOrders(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			actual, err := s.service.GetUserOrders(tc.userID)
			if tc.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().NotNil(actual)
			}
		})
	}
}

func (s *OrderTestSuite) TestEnqueueOrdersFromAccrualSystem() {
	type testCase struct {
		name         string
		mockBehavior func()
		wantErr      bool
	}
	testCases := []testCase{
		{
			name: "OK",
			mockBehavior: func() {
				s.repository.EXPECT().
					FindManyByStatus(model.StatusNew, model.StatusProcessing).
					Return(&[]model.Order{
						{ID: "79927398713", UserID: 1, OrderStatus: model.StatusNew},
						{ID: "79927398714", UserID: 2, OrderStatus: model.StatusProcessing},
					}, nil)
				s.semaphore.EXPECT().Acquire().AnyTimes()
				s.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				s.restClient.EXPECT().Get(gomock.Any()).AnyTimes().
					DoAndReturn(func(matcher string) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte(`{"order":"79927398713","status":"PROCESSED","accrual":150.0}`))),
						}, nil
					})
				s.semaphore.EXPECT().Release().MinTimes(1)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mockBehavior: func() {
				s.repository.EXPECT().
					FindManyByStatus(model.StatusNew, model.StatusProcessing).
					Return(nil, errors.New("some error"))
				s.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				s.logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			},
			wantErr: true,
		},
		{
			name: "no orders to process",
			mockBehavior: func() {
				s.repository.EXPECT().
					FindManyByStatus(model.StatusNew, model.StatusProcessing).
					Return(&[]model.Order{}, nil)
				s.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			err := s.service.EnqueueOrdersFromAccrualSystem(context.TODO(), &sync.WaitGroup{})
			// Wait for all goroutines to finish
			time.Sleep(200 * time.Millisecond)
			if tc.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}
