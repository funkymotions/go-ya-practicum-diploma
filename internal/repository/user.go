package repository

import (
	"database/sql"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type userRepository struct {
	driver SQLExecutor
}

func NewUserRepository(d SQLExecutor) *userRepository {
	return &userRepository{
		driver: d,
	}
}

func (r *userRepository) FindOneByID(userID uint) (*model.User, error) {
	sqlString := "SELECT id, login, password_hash FROM users WHERE id = $1"
	row := r.driver.QueryRow(sqlString, userID)
	result := model.User{}
	err := row.Scan(&result.ID, &result.Login, &result.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *userRepository) FindByLogin(login string) (*model.User, error) {
	sqlString := "SELECT id, login, password_hash FROM users WHERE login = $1"
	row := r.driver.QueryRow(sqlString, login)
	result := model.User{}
	err := row.Scan(&result.ID, &result.Login, &result.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *userRepository) Create(login, passwordHash string) (*model.User, error) {
	sqlString := "INSERT INTO users (login, password_hash) VALUES ($1, $2) ON CONFLICT (login) DO NOTHING RETURNING id, login"
	row := r.driver.QueryRow(sqlString, login, passwordHash)
	result := model.User{}
	err := row.Scan(&result.ID, &result.Login)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}
