package service

import (
	"errors"
	"testing"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type accountServiceTestSuite struct {
	suite.Suite
	service        *accountService
	orderRepo      *mocks.MockOrderRepository
	withdrawalRepo *mocks.MockWithdrawalRepository
}

func (s *accountServiceTestSuite) SetupTest() {
	s.orderRepo = mocks.NewMockOrderRepository(gomock.NewController(s.T()))
	s.withdrawalRepo = mocks.NewMockWithdrawalRepository(gomock.NewController(s.T()))
	s.service = NewAccountService(s.orderRepo, s.withdrawalRepo)
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(accountServiceTestSuite))
}

func (s *accountServiceTestSuite) TestCalculateAccountBalance() {
	type testCase struct {
		name         string
		userID       uint
		mockBehavior func()
		expectedSum  float64
		wantErr      bool
	}
	testCases := []testCase{
		{
			name:   "no orders or withdrawals",
			userID: 1,
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(0.0, nil)
				s.withdrawalRepo.EXPECT().
					GetUserWithdrawalBalance(gomock.Any()).
					Return(0.0, nil)
			},
			expectedSum: 0.0,
			wantErr:     false,
		},
		{
			name:   "orders and withdrawals present",
			userID: 2,
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(150.0, nil)
				s.withdrawalRepo.EXPECT().
					GetUserWithdrawalBalance(gomock.Any()).
					Return(50.0, nil)
			},
			expectedSum: 100.0,
			wantErr:     false,
		},
		{
			name:   "withdrawals repository error",
			userID: 1,
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(0.0, nil)
				s.withdrawalRepo.EXPECT().
					GetUserWithdrawalBalance(gomock.Any()).
					Return(0.0, errors.New("some error"))
			},
			expectedSum: 0.0,
			wantErr:     true,
		},
		{
			name:   "order repository error",
			userID: 1,
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(0.0, errors.New("some error"))
			},
			expectedSum: 0.0,
			wantErr:     true,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			balance, _, err := s.service.CalculateAccountBalance(tc.userID)
			if tc.wantErr {
				s.Assert().NotNil(err)
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().Equal(tc.expectedSum, balance)
			}
		})
	}
}

func (s *accountServiceTestSuite) TestWithdrawFromAccount() {
	type testCase struct {
		name         string
		userID       uint
		parameters   *dto.AccountBalanceWithdrawRequest
		mockBehavior func()
		wantErr      bool
	}
	testCases := []testCase{
		{
			name:   "sufficient funds",
			userID: 1,
			parameters: &dto.AccountBalanceWithdrawRequest{
				Order: "79927398713",
				Sum:   50.0,
			},
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(100.0, nil)
				s.withdrawalRepo.EXPECT().
					GetUserWithdrawalBalance(gomock.Any()).
					Return(0.0, nil)
				s.withdrawalRepo.EXPECT().
					CreateWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "insufficient funds",
			userID: 1,
			parameters: &dto.AccountBalanceWithdrawRequest{
				Order: "79927398713",
				Sum:   150.0,
			},
			mockBehavior: func() {
				s.orderRepo.EXPECT().
					GetUserOrderBalance(gomock.Any()).
					Return(100.0, nil)
				s.withdrawalRepo.EXPECT().
					GetUserWithdrawalBalance(gomock.Any()).
					Return(0.0, nil)
			},
			wantErr: true,
		},
		{
			name:   "invalid order ID",
			userID: 1,
			parameters: &dto.AccountBalanceWithdrawRequest{
				Order: "799273987130",
				Sum:   50.0,
			},
			mockBehavior: nil,
			wantErr:      true,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			err := s.service.WithdrawFromAccount(tc.userID, tc.parameters)
			if tc.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}
