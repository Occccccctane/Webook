package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Dao"
	"context"
)

var (
	ErrUserNotFound = Dao.ErrRecordNotFound
	EmailUniqueErr  = Dao.EmailUniqueErr
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

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (Domain.User, error) {
	user, err1 := repo.dao.EmailSearch(ctx, email)
	if err1 != nil {
		return Domain.User{}, err1
	}
	return repo.toDomain(user), nil
}

func (repo *UserRepository) Edit(ctx context.Context, u Domain.User) error {
	//查找要修改的记录
	user, err1 := repo.dao.EmailSearch(ctx, u.Email)
	if err1 != nil {
		return err1
	}

	//更新信息
	err2 := repo.dao.Update(Dao.User{
		Id:       user.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Info:     u.Info,
		Ctime:    user.Ctime,
		Utime:    0,
	})
	if err2 != nil {
		return err2
	}
	return nil
}

// 将dao的实体转换成domain的实体，避免跨层调用
func (repo *UserRepository) toDomain(user Dao.User) Domain.User {
	return Domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
}
