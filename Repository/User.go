package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Dao"
	"context"
)

type UserRepository struct {
	dao *Dao.UserDao
}

func NewUserRepository(dao *Dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, Dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
