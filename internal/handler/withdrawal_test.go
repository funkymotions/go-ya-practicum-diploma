package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler/mocks"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type WithdrawalHandlerTestSuite struct {
	suite.Suite
	handler *withdrawalHandler
	service *mocks.MockWithdrawalService
}

func (s *WithdrawalHandlerTestSuite) SetupTest() {
	s.service = mocks.NewMockWithdrawalService(gomock.NewController(s.T()))
	s.handler = NewWithdrawalHandler(s.service)
}

func TestWithdrawalHandler(t *testing.T) {
	suite.Run(t, new(WithdrawalHandlerTestSuite))
}

func (s *WithdrawalHandlerTestSuite) TestGetWithdrawals() {
	type testCase struct {
		name         string
		userID       string
		contentType  string
		wantStatus   int
		mockBehavior func()
	}
	testCases := []testCase{
		{
			name:       "no withdrawals",
			userID:     "1",
			wantStatus: http.StatusNoContent,
			mockBehavior: func() {
				s.service.EXPECT().
					GetUserWithdrawals(uint(1)).
					Return(nil, apperrors.ErrNoWithdrawals)
			},
		},
		{
			name:       "successful retrieval",
			userID:     "1",
			wantStatus: http.StatusOK,
			mockBehavior: func() {
				curTime := time.Now()
				s.service.EXPECT().
					GetUserWithdrawals(uint(1)).
					Return(&[]model.Withdrawal{
						{
							OrderID:     "79927398713",
							Amount:      100.0,
							ProcessedAt: &curTime,
						},
					}, nil)
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior()
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, middleware.UserIDValue("userID"), tc.userID)
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/withdrawals", nil)
			w := httptest.NewRecorder()
			s.handler.GetUserWithdrawals(w, req)
			resp := w.Result()
			s.Assert().Equal(tc.wantStatus, resp.StatusCode)
			resp.Body.Close()
		})
	}
}
