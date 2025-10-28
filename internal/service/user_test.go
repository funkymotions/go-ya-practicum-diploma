package service

import (
	"errors"
	"testing"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service/mocks"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type userServiceTestSuite struct {
	suite.Suite
	repository *mocks.MockUserRepository
	service    *userService
}

func (s *userServiceTestSuite) SetupTest() {
	s.repository = mocks.NewMockUserRepository(
		gomock.NewController(s.T()),
	)
	s.service = NewUserService(s.repository)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(userServiceTestSuite))
}

func (s *userServiceTestSuite) TestRegister() {
	type testCase struct {
		name        string
		args        *dto.RegisterUserRequest
		mockExpects func()
		wantErr     bool
	}
	testCases := []testCase{
		{
			name: "OK",
			args: &dto.RegisterUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					Create("user1", gomock.Any()).
					Return(&model.User{ID: 1, Login: "user1"}, nil)
			},
			wantErr: false,
		},
		{
			name: "Repository error",
			args: &dto.RegisterUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					Create("user1", gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockExpects != nil {
				tc.mockExpects()
			}
			actual, err := s.service.Register(tc.args)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tc.args.Login, actual.Login)
			}
		})
	}
}

func (s *userServiceTestSuite) TestLogin() {
	type testCase struct {
		name        string
		args        *dto.LoginUserRequest
		mockExpects func()
		wantErr     bool
	}
	testCases := []testCase{
		{
			name: "OK",
			args: &dto.LoginUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					FindByLogin("user1").
					Return(
						&model.User{ID: 1, Login: "user1", PasswordHash: utils.HashPassword("pass123")},
						nil,
					)
			},
			wantErr: false,
		},
		{
			name: "Unknown repository error",
			args: &dto.LoginUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					FindByLogin("user1").
					Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name: "No user found",
			args: &dto.LoginUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					FindByLogin("user1").
					Return(nil, apperrors.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "Wrong password",
			args: &dto.LoginUserRequest{
				Login:    "user1",
				Password: "wrongpass",
			},
			mockExpects: func() {
				s.repository.EXPECT().
					FindByLogin("user1").
					Return(
						&model.User{ID: 1, Login: "user1", PasswordHash: utils.HashPassword("pass123")},
						nil,
					)
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.mockExpects != nil {
				tc.mockExpects()
			}
			actual, err := s.service.Login(tc.args)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tc.args.Login, actual.Login)
			}
		})
	}
}
