package service

import (
	"errors"
	"testing"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type withdrawalServiceTestSuite struct {
	suite.Suite
	repository *mocks.MockWithdrawalRepository
	service    *withdrawalService
}

func (s *withdrawalServiceTestSuite) SetupTest() {
	s.repository = mocks.NewMockWithdrawalRepository(gomock.NewController(s.T()))
	s.service = NewWithdrawalService(s.repository)
}

func TestWithdrawalServiceTestSuite(t *testing.T) {
	suite.Run(t, new(withdrawalServiceTestSuite))
}

func (s *withdrawalServiceTestSuite) TestGetUserWithdrawals() {
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
				serviceReturn := &[]model.Withdrawal{
					{
						ID:      1,
						UserID:  1,
						OrderID: "79927398713",
						Amount:  100.0,
					},
				}
				s.repository.EXPECT().
					GetUserWithdrawals(gomock.Any()).
					Return(serviceReturn, nil)
			},
			wantErr: false,
		},
		{
			name:   "no withdrawals",
			userID: 1,
			mockBehavior: func() {
				serviceReturn := &[]model.Withdrawal{}
				s.repository.EXPECT().
					GetUserWithdrawals(gomock.Any()).
					Return(serviceReturn, nil)
			},
			wantErr: true,
		},
		{
			name:   "service unexpected error",
			userID: 1,
			mockBehavior: func() {
				s.repository.EXPECT().
					GetUserWithdrawals(gomock.Any()).
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
			_, err := s.service.GetUserWithdrawals(tc.userID)
			if tc.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
