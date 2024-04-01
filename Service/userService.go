package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"context"
	"golang.org/x/crypto/bcrypt"
)

var EmailUniqueErr = Repository.EmailUniqueErr

type UserService struct {
	repo *Repository.UserRepository
}

func NewUserService(repo *Repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u Domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//没报错就将密码加密为哈希，将哈希存入数据库中
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
