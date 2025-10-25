package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler/mocks"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type accountHandlerTestSuite struct {
	suite.Suite
	service *mocks.MockaccountService
	handler *accountHandler
}

func (s *accountHandlerTestSuite) SetupTest() {
	s.service = mocks.NewMockaccountService(gomock.NewController(s.T()))
	s.handler = NewAccountHandler(s.service)
}

func TestAccountHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(accountHandlerTestSuite))
}

func (s *accountHandlerTestSuite) TestAccountBalance() {
	type request struct {
		userID string
	}
	type testCase struct {
		name               string
		request            request
		expectedStatusCode int
		mockBehavior       func(r *request)
	}
	testCases := []testCase{
		{
			name: "OK",
			request: request{
				userID: "1",
			},
			expectedStatusCode: http.StatusOK,
			mockBehavior: func(r *request) {
				s.service.
					EXPECT().
					CalculateAccountBalance(uint(1)).
					Return(100.0, 100.0, nil)
			},
		},
		{
			name: "Wrong userID in context",
			request: request{
				userID: "***",
			},
			expectedStatusCode: http.StatusUnauthorized,
			mockBehavior:       nil,
		},
		{
			name: "service error",
			request: request{
				userID: "1",
			},
			expectedStatusCode: http.StatusInternalServerError,
			mockBehavior: func(r *request) {
				s.service.
					EXPECT().
					CalculateAccountBalance(uint(1)).
					Return(0.0, 0.0, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior(&tc.request)
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, middleware.UserIDValue("userID"), tc.request.userID)
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/balance", nil)
			w := httptest.NewRecorder()
			s.handler.GetAccountBalance(w, req)
			s.Assert().Equal(tc.expectedStatusCode, w.Result().StatusCode)
			req.Body.Close()
		})
	}
}

func (s *accountHandlerTestSuite) TestAccountWithdrawal() {
	type request struct {
		userID      string
		body        string
		contentType string
	}
	type testCase struct {
		name               string
		request            request
		expectedStatusCode int
		mockBehavior       func(r *request)
	}

	testCases := []testCase{
		{
			name: "OK",
			request: request{
				userID:      "1",
				body:        `{"order":"79927398713","sum":50.0}`,
				contentType: "application/json",
			},
			expectedStatusCode: http.StatusOK,
			mockBehavior: func(r *request) {
				s.service.
					EXPECT().
					WithdrawFromAccount(uint(1), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "Invalid request body",
			request: request{
				userID:      "1",
				body:        `{"order":"123","sum":0.0}`,
				contentType: "application/json",
			},
			expectedStatusCode: http.StatusBadRequest,
			mockBehavior:       nil,
		},
		{
			name: "Invalid content type",
			request: request{
				userID:      "1",
				body:        `{"order":"123","sum":0.0}`,
				contentType: "application/xml",
			},
			expectedStatusCode: http.StatusBadRequest,
			mockBehavior:       nil,
		},
		{
			name: "Invalid userID",
			request: request{
				userID:      "***",
				body:        ``,
				contentType: "application/json",
			},
			expectedStatusCode: http.StatusUnauthorized,
			mockBehavior:       nil,
		},
		{
			name: "Wrong sum to withdraw",
			request: request{
				userID:      "1",
				body:        `{"order":"79927398713","sum":50.0}`,
				contentType: "application/json",
			},
			expectedStatusCode: apperrors.ErrAccountInsufficientFunds.StatusCode,
			mockBehavior: func(r *request) {
				s.service.
					EXPECT().
					WithdrawFromAccount(uint(1), gomock.Any()).
					Return(apperrors.ErrAccountInsufficientFunds)
			},
		},
		{
			name: "unexpected service error",
			request: request{
				userID:      "1",
				body:        `{"order":"79927398713","sum":50.0}`,
				contentType: "application/json",
			},
			expectedStatusCode: http.StatusInternalServerError,
			mockBehavior: func(r *request) {
				s.service.
					EXPECT().
					WithdrawFromAccount(uint(1), gomock.Any()).
					Return(errors.New("some unexpected error"))
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior(&tc.request)
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, middleware.UserIDValue("userID"), tc.request.userID)
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(tc.request.body))
			req.Header.Set("Content-Type", tc.request.contentType)
			w := httptest.NewRecorder()
			s.handler.WithdrawAccount(w, req)
			s.Assert().Equal(tc.expectedStatusCode, w.Result().StatusCode)
		})
	}
}
