package interfaces

import (
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type UserService interface {
	Register(*dto.RegisterUserRequest) (*model.User, error)
	Login(*dto.LoginUserRequest) (*model.User, error)
}

type UserRepository interface {
	FindOneByID(id uint) (*model.User, error)
	Create(string, string) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
}
