package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Cache"
	"GinStart/Repository/Dao"
	"context"
	"log"
)

var (
	ErrUserNotFound = Dao.ErrRecordNotFound
	EmailUniqueErr  = Dao.EmailUniqueErr
)

type UserRepository struct {
	dao   *Dao.UserDao
	cache *Cache.UserCache
}

// NewUserRepository 类似构造函数，将数据暴露出来
func NewUserRepository(dao *Dao.UserDao, cache *Cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, Dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (Domain.User, error) {
	user, err1 := repo.dao.FindByEmail(ctx, email)
	if err1 != nil {
		return Domain.User{}, err1
	}
	return repo.toDomain(user), nil
}

func (repo *UserRepository) FindByID(ctx context.Context, uid int64) (Domain.User, error) {
	// 从缓存中读
	u, err := repo.cache.Get(ctx, uid)
	if err != nil {
		return u, err
	}
	// 缓存中没有，从数据库中读
	user, err1 := repo.dao.FindByID(ctx, uid)
	if err1 != nil {
		return Domain.User{}, err1
	}
	// 从数据库中读出来后，转换格式放入缓存中
	u = repo.toDomain(user)
	err2 := repo.cache.Set(ctx, u)
	//将错误忽略
	if err2 != nil {
		log.Println(err)
	}
	return u, nil
}

// 避免缓存击穿的写法
func (repo *UserRepository) FindByIDV1(ctx context.Context, uid int64) (Domain.User, error) {
	// 从缓存中读
	u, err := repo.cache.Get(ctx, uid)
	switch err {
	case nil:
		return u, err
	case Cache.ErrKeyNotExist:
		// 缓存中没有，从数据库中读
		user, err1 := repo.dao.FindByID(ctx, uid)
		if err1 != nil {
			return Domain.User{}, err1
		}
		// 从数据库中读出来后，转换格式放入缓存中
		u = repo.toDomain(user)
		err2 := repo.cache.Set(ctx, u)
		//将错误忽略
		if err2 != nil {
			log.Println(err)
		}
		return u, nil
	default:
		//不知道Redis有没有数据，但是知道Redis不正常
		//接近于降级的写法
		return Domain.User{}, err

	}

}
func (repo *UserRepository) Edit(ctx context.Context, u Domain.User) error {
	//查找要修改的记录
	user, err1 := repo.dao.FindByID(ctx, u.Id)
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
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		Info:     user.Info,
	}
}
