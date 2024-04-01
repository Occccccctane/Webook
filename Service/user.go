package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"context"
)

type UserService struct {
	repo *Repository.UserRepository
}

func NewUserService(repo *Repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u Domain.User) error {
	return svc.repo.Create(ctx, u)
}
