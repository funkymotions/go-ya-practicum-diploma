package service

import (
	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
)

type userRepository interface {
	FindOneByID(id uint) (*model.User, error)
	Create(string, string) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
}

type userService struct {
	repo userRepository
}

func NewUserService(r userRepository) *userService {
	return &userService{
		repo: r,
	}
}

func (s *userService) Register(req *dto.RegisterUserRequest) (*model.User, error) {
	hashed := utils.HashPassword(req.Password)
	user, err := s.repo.Create(req.Login, hashed)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(req *dto.LoginUserRequest) (*model.User, error) {
	user, err := s.repo.FindByLogin(req.Login)
	if err == apperrors.ErrUserNotFound {
		return nil, apperrors.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !utils.ComparePasswords(user.PasswordHash, req.Password) {
		return nil, apperrors.ErrInvalidCredentials
	}
	return user, nil
}
