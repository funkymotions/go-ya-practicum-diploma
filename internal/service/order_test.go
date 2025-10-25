package service

import (
	"errors"
	"testing"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderTestSuite struct {
	suite.Suite
	repository *mocks.MockorderRepository
	service    *orderService
}

func (s *OrderTestSuite) SetupTest() {
	s.repository = mocks.NewMockorderRepository(gomock.NewController(s.T()))
	s.service = NewOrderService(s.repository)
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
