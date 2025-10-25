package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler/mocks"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type userHandlerTestSuite struct {
	suite.Suite
	service *mocks.MockuserService
	server  *httptest.Server
	handler *userHandler
}

func (s *userHandlerTestSuite) SetupTest() {
	s.service = mocks.NewMockuserService(
		gomock.NewController(
			s.T(),
		),
	)
	s.handler = NewUserHandler(s.service)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(userHandlerTestSuite))
}

func (s *userHandlerTestSuite) TestRegisterHandler() {
	type testCase struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		mockExpects func()
	}
	testCases := []testCase{
		{
			name:        "valid request",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusOK,
			mockExpects: func() {
				s.service.EXPECT().
					Register(&dto.RegisterUserRequest{
						Login:    "user1",
						Password: "pass123",
					}).
					Return(&model.User{ID: 1, Login: "user1"}, nil)
			},
		},
		{
			name:        "service error: user already exists",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusConflict,
			mockExpects: func() {
				s.service.EXPECT().
					Register(&dto.RegisterUserRequest{
						Login:    "user1",
						Password: "pass123",
					}).
					Return(nil, apperrors.ErrUserAlreadyExists)
			},
		},
		{
			name:        "service error: unknown errror",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusInternalServerError,
			mockExpects: func() {
				s.service.EXPECT().
					Register(&dto.RegisterUserRequest{
						Login:    "user1",
						Password: "pass123",
					}).
					Return(nil, errors.New("some unknown error"))
			},
		},
		{
			name:        "invalid json body",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"{{{}}}`,
			wantStatus:  http.StatusBadRequest,
			mockExpects: nil,
		},
		{
			name:        "invalid content type",
			contentType: "text/html",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusBadRequest,
			mockExpects: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockExpects != nil {
				tc.mockExpects()
			}
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			rr := httptest.NewRecorder()
			s.handler.Register(rr, req)
			s.Assert().Equal(tc.wantStatus, rr.Code)
		})
	}
}

func (s *userHandlerTestSuite) TestLoginHandler() {
	type testCase struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		mockExpects func()
	}
	testCases := []testCase{
		{
			name:        "valid request",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusOK,
			mockExpects: func() {
				s.service.EXPECT().
					Login(&dto.LoginUserRequest{
						Login:    "user1",
						Password: "pass123",
					}).
					Return(&model.User{ID: 1, Login: "user1"}, nil)
			},
		},
		{
			name:        "service error: invalid credentials",
			contentType: "application/json",
			body:        `{"login":"user1","password":"wrongpass"}`,
			wantStatus:  http.StatusUnauthorized,
			mockExpects: func() {
				s.service.EXPECT().
					Login(&dto.LoginUserRequest{
						Login:    "user1",
						Password: "wrongpass",
					}).
					Return(nil, apperrors.ErrInvalidCredentials)
			},
		},
		{
			name:        "service error: unknown errror",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"}`,
			wantStatus:  http.StatusInternalServerError,
			mockExpects: func() {
				s.service.EXPECT().
					Login(&dto.LoginUserRequest{
						Login:    "user1",
						Password: "pass123",
					}).
					Return(nil, errors.New("some unknown error"))
			},
		},
		{
			name:        "invalid json body",
			contentType: "application/json",
			body:        `{"login":"user1","password":"pass123"{{{}}}`,
			wantStatus:  http.StatusBadRequest,
			mockExpects: nil,
		},
		{
			name:        "invalid content type",
			contentType: "text/html",
			body:        `{"login":"user1","password":"pass123"}"`,
			wantStatus:  http.StatusBadRequest,
			mockExpects: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockExpects != nil {
				tc.mockExpects()
			}
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			rr := httptest.NewRecorder()
			s.handler.Login(rr, req)
			s.Assert().Equal(tc.wantStatus, rr.Code)
		})
	}
}
