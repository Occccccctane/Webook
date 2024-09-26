package Service

import (
	"GinStart/Domain"
	"GinStart/Repository"
	"GinStart/pkg/logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserUnique            = Repository.ErrUserUnique
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
)

type UserService interface {
	Signup(ctx context.Context, u Domain.User) error
	Login(context context.Context, email, password string) (Domain.User, error)
	Edit(ctx context.Context, newPassword string, u Domain.User) error
	FindOrCreate(c *gin.Context, phone string) (Domain.User, error)
	FindById(ctx *gin.Context, id int64) (Domain.User, error)
}

type userService struct {
	repo   Repository.UserRepository
	logger logger.Logger
}

func NewUserService(repo Repository.UserRepository, l logger.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: l,
	}
}
func (svc *userService) FindById(ctx *gin.Context, uid int64) (Domain.User, error) {
	user, err := svc.repo.FindByID(ctx, uid)
	if err != nil {
		return Domain.User{}, err
	}
	return user, nil
}

func (svc *userService) Signup(ctx context.Context, u Domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//没报错就将密码加密为哈希，将哈希存入数据库中
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(context context.Context, email, password string) (Domain.User, error) {
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

func (svc *userService) Edit(ctx context.Context, newPassword string, u Domain.User) error {
	//验证原始密码
	user, err1 := svc.repo.FindByEmail(ctx, u.Email)
	if err1 == Repository.ErrUserNotFound {
		return ErrInvalidUserOrPassword
	}
	if err1 != nil {
		return err1
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err2 != nil {
		return ErrInvalidUserOrPassword
	}
	//加密新密码
	hash, err3 := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err3 != nil {
		return err3
	}

	//修改信息
	//如果有新密码就保存新哈希，将哈希存入数据库中
	if newPassword != u.Password {
		u.Password = string(hash)
	}

	err4 := svc.repo.Edit(ctx, u)
	if err4 != nil {
		return err4
	}
	return nil
}

func (svc *userService) FindOrCreate(c *gin.Context, phone string) (Domain.User, error) {
	//先找一下，我们认为，大部分用户是已经存在的
	u, err := svc.repo.FindByPhone(c, phone)
	if err != Repository.ErrUserNotFound {
		return u, err
	}
	//找不到意味着是一个新用户
	//svc.logger.Info("新用户", )
	err = svc.repo.Create(c, Domain.User{
		Phone: phone,
	})
	//有err但不是唯一索引冲突
	//系统错误
	if err != nil && err != Repository.ErrUserUnique {
		return Domain.User{}, err
	}
	//要么没错误，要么唯一索引限制即用户存在
	//可能有主从延迟，理论上讲强行走主库
	return svc.repo.FindByPhone(c, phone)
}
