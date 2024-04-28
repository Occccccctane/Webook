package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Cache"
	"GinStart/Repository/Dao"
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	ErrUserNotFound = Dao.ErrRecordNotFound
	ErrUserUnique   = Dao.EmailUniqueErr
)

type UserRepository interface {
	Create(ctx context.Context, u Domain.User) error
	FindByEmail(ctx context.Context, email string) (Domain.User, error)
	FindByID(ctx context.Context, uid int64) (Domain.User, error)
	FindByPhone(ctx context.Context, phone string) (Domain.User, error)
	Edit(ctx context.Context, u Domain.User) error
}
type CacheUserRepository struct {
	dao   Dao.UserDao
	cache Cache.UserCache
}

// NewCacheUserRepository 类似构造函数，将数据暴露出来
func NewCacheUserRepository(dao Dao.UserDao, cache Cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (Domain.User, error) {
	user, err1 := repo.dao.FindByEmail(ctx, email)
	if err1 != nil {
		return Domain.User{}, err1
	}
	return repo.toDomain(user), nil
}

func (repo *CacheUserRepository) FindByID(ctx context.Context, uid int64) (Domain.User, error) {
	// 从缓存中读
	u, err := repo.cache.Get(ctx, uid)
	if err == nil {
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
func (repo *CacheUserRepository) FindByIDV1(ctx context.Context, uid int64) (Domain.User, error) {
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

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (Domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return Domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CacheUserRepository) Edit(ctx context.Context, u Domain.User) error {
	//查找要修改的记录
	user, err1 := repo.dao.FindByID(ctx, u.Id)
	if err1 != nil {
		return err1
	}
	u1 := repo.toEntity(u)
	u1.Ctime = user.Ctime

	//更新信息
	err2 := repo.dao.Update(u1)
	if err2 != nil {
		return err2
	}
	return nil
}

// 将dao的实体转换成domain的实体，避免跨层调用
func (repo *CacheUserRepository) toDomain(user Dao.User) Domain.User {
	return Domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		Info:     user.Info,
		Phone:    user.Phone.String,
	}
}

func (repo *CacheUserRepository) toEntity(u Domain.User) Dao.User {
	return Dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Info:     u.Info,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Utime: time.Now().UnixMilli(),
	}
}
