package repository

import "github.com/funkymotions/go-ya-practicum-diploma/internal/model"

type accountRepository struct {
	driver SQLExecutor
}

func NewAccountRepository(d SQLExecutor) *accountRepository {
	return &accountRepository{
		driver: d,
	}
}

func (r *accountRepository) FindOneByID(accountID uint) (*model.Account, error) {
	panic("implement me")
}
