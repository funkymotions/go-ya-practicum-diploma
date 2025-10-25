package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler/mocks"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type orderHandlerTestSuite struct {
	suite.Suite
	handler *orderHandler
	service *mocks.MockorderService
}

func (s *orderHandlerTestSuite) SetupTest() {
	s.service = mocks.NewMockorderService(gomock.NewController(s.T()))
	s.handler = NewOrderHandler(s.service)
}

func TestOrderHandler(t *testing.T) {
	suite.Run(t, new(orderHandlerTestSuite))
}

func (s *orderHandlerTestSuite) TestOrderRegistration() {
	type requestArgs struct {
		body               io.Reader
		requestContentType string
		userID             uint
	}
	type testCase struct {
		name           string
		wantStatusCode int
		serviceArg     *dto.RegisterOrderRequest
		mockBehavior   func(*dto.RegisterOrderRequest, *requestArgs)
		requestArgs    requestArgs
	}

	testCases := []testCase{
		{
			name: "successful registration",
			requestArgs: requestArgs{
				requestContentType: "text/plain",
				userID:             1,
				body:               io.NopCloser(strings.NewReader("79927398713")),
			},
			wantStatusCode: http.StatusAccepted,
			serviceArg:     &dto.RegisterOrderRequest{},
			mockBehavior: func(serviceArg *dto.RegisterOrderRequest, req *requestArgs) {
				// preconditions
				data, _ := io.ReadAll(req.body)
				serviceArg.OrderID = string(data)
				serviceArg.UserID = req.userID
				req.body = io.NopCloser(strings.NewReader(string(data)))

				// mock expectations
				s.service.EXPECT().
					RegisterOrder(serviceArg).
					Return(nil)
			},
		},
		{
			name: "error: wrong order number",
			requestArgs: requestArgs{
				requestContentType: "text/plain",
				userID:             1,
				body:               io.NopCloser(strings.NewReader("799273987130")),
			},
			wantStatusCode: http.StatusUnprocessableEntity,
			serviceArg:     &dto.RegisterOrderRequest{},
			mockBehavior: func(serviceArg *dto.RegisterOrderRequest, req *requestArgs) {
				// preconditions
				data, _ := io.ReadAll(req.body)
				serviceArg.OrderID = string(data)
				serviceArg.UserID = req.userID
				req.body = io.NopCloser(strings.NewReader(string(data)))

				// mock expectations
				s.service.EXPECT().
					RegisterOrder(serviceArg).
					Return(apperrors.ErrOrderInvalidID)
			},
		},
		{
			name: "error: unexpected error",
			requestArgs: requestArgs{
				requestContentType: "text/plain",
				userID:             1,
				body:               io.NopCloser(strings.NewReader("79927398713")),
			},
			wantStatusCode: http.StatusInternalServerError,
			serviceArg:     &dto.RegisterOrderRequest{},
			mockBehavior: func(serviceArg *dto.RegisterOrderRequest, req *requestArgs) {
				// preconditions
				data, _ := io.ReadAll(req.body)
				serviceArg.OrderID = string(data)
				serviceArg.UserID = req.userID
				req.body = io.NopCloser(strings.NewReader(string(data)))

				// mock expectations
				s.service.EXPECT().
					RegisterOrder(serviceArg).
					Return(errors.New("some unexpected error"))
			},
		},
		{
			name: "wrong content type",
			requestArgs: requestArgs{
				requestContentType: "application/json",
				userID:             1,
				body:               io.NopCloser(strings.NewReader("79927398713")),
			},
			wantStatusCode: http.StatusBadRequest,
			serviceArg:     &dto.RegisterOrderRequest{},
			mockBehavior:   nil,
		},
		{
			name: "empty body",
			requestArgs: requestArgs{
				requestContentType: "text/plain",
				userID:             1,
				body:               nil,
			},
			wantStatusCode: http.StatusBadRequest,
			serviceArg:     &dto.RegisterOrderRequest{},
			mockBehavior:   nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior(tc.serviceArg, &tc.requestArgs)
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, middleware.UserIDValue("userID"), fmt.Sprintf("%d", tc.requestArgs.userID))
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/user/orders", tc.requestArgs.body)
			req.Header.Set("Content-Type", tc.requestArgs.requestContentType)
			w := httptest.NewRecorder()
			s.handler.RegisterOrder(w, req)
			s.Assert().Equal(tc.wantStatusCode, w.Result().StatusCode)
			req.Body.Close()
		})
	}
}

func (s *orderHandlerTestSuite) TestGetUserOrders() {
	type requestArgs struct {
		userID uint
	}
	type testCase struct {
		name           string
		wantStatusCode int
		mockBehavior   func(*requestArgs)
		requestArgs    requestArgs
	}

	testCases := []testCase{
		{
			name: "successful retrieval",
			requestArgs: requestArgs{
				userID: 1,
			},
			wantStatusCode: http.StatusOK,
			mockBehavior: func(req *requestArgs) {
				s.service.EXPECT().
					GetUserOrders(req.userID).
					Return(&[]model.Order{}, nil)
			},
		},
		{
			name: "error: no orders",
			requestArgs: requestArgs{
				userID: 1,
			},
			wantStatusCode: http.StatusNoContent,
			mockBehavior: func(req *requestArgs) {
				s.service.EXPECT().
					GetUserOrders(req.userID).
					Return(nil, apperrors.ErrOrderListEmpty)
			},
		},
		{
			name: "error: unexpected error",
			requestArgs: requestArgs{
				userID: 1,
			},
			wantStatusCode: http.StatusInternalServerError,
			mockBehavior: func(req *requestArgs) {
				s.service.EXPECT().
					GetUserOrders(req.userID).
					Return(nil, errors.New("some unexpected error"))
			},
		},
		{
			name: "OK with order list",
			requestArgs: requestArgs{
				userID: 1,
			},
			wantStatusCode: http.StatusOK,
			mockBehavior: func(req *requestArgs) {
				accrual := float64(100.0)
				ret := &[]model.Order{
					{
						ID:          "79927398713",
						UserID:      req.userID,
						OrderStatus: model.StatusNew,
						CreatedAt:   &time.Time{},
						Accrual:     &accrual,
					},
				}
				s.service.EXPECT().
					GetUserOrders(req.userID).
					Return(ret, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockBehavior != nil {
				tc.mockBehavior(&tc.requestArgs)
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, middleware.UserIDValue("userID"), fmt.Sprintf("%d", tc.requestArgs.userID))
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/user/orders", nil)
			w := httptest.NewRecorder()
			s.handler.GetUserOrders(w, req)
			s.Assert().Equal(tc.wantStatusCode, w.Result().StatusCode)
			req.Body.Close()
		})
	}
}
