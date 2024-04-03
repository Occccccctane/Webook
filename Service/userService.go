package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	EmailUniqueErr           = Repository.EmailUniqueErr
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
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
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//没报错就将密码加密为哈希，将哈希存入数据库中
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(context context.Context, email, password string) (Domain.User, error) {
	user, err1 := svc.repo.FindByEmail(context, email)
	if err1 == Repository.ErrUserNotFound {
		return Domain.User{}, ErrInvalidUserOrPassword
	}
	if err1 != nil {
		return Domain.User{}, err1
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		return Domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *UserService) Edit(ctx context.Context, email, password, newPassword string) (Domain.User, error) {
	//验证原始密码
	user, err1 := svc.repo.FindByEmail(ctx, email)
	if err1 == Repository.ErrUserNotFound {
		return Domain.User{}, ErrInvalidUserOrPassword
	}
	if err1 != nil {
		return Domain.User{}, err1
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		return Domain.User{}, ErrInvalidUserOrPassword
	}

	hash, err3 := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err3 != nil {
		return Domain.User{}, err3
	}

	//修改信息
	//没报错就将密码加密为哈希，将哈希存入数据库中
	newPassword = string(hash)
	newUser, err4 := svc.repo.Edit(ctx, email, newPassword)
	if err4 != nil {
		return Domain.User{}, err4
	}
	return newUser, nil
}
