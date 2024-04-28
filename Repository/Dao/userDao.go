package Dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	EmailUniqueErr    = errors.New("邮箱唯一错误")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(context context.Context, email string) (User, error)
	Update(user User) error
	FindByID(ctx context.Context, uid int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{db: db}
}

func (dao *GormUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error

	//将错误类型断言为unique键错误，并给出特定的错误处理
	if mysl, ok := err.(*mysql.MySQLError); ok {
		const uniqueErrNum uint16 = 1062
		if mysl.Number == uniqueErrNum {
			return EmailUniqueErr
		}
	}
	return err
}

func (dao *GormUserDao) FindByEmail(context context.Context, email string) (User, error) {
	var user User
	err1 := dao.db.WithContext(context).Where("email=?", email).First(&user).Error
	if err1 != nil {
		return user, err1
	}
	return user, nil
}

func (dao *GormUserDao) Update(user User) error {
	//更新信息

	now := time.Now().UnixMilli()
	user.Utime = now

	err1 := dao.db.Save(&user).Error
	if err1 != nil {
		return err1
	}
	return nil
}

func (dao *GormUserDao) FindByID(ctx context.Context, uid int64) (User, error) {
	var user User
	err1 := dao.db.WithContext(ctx).Where("ID=?", uid).First(&user).Error
	if err1 != nil {
		return user, err1
	}
	return user, nil
}

func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err1 := dao.db.WithContext(ctx).Where("phone=?", phone).First(&user).Error
	if err1 != nil {
		return user, err1
	}
	return user, nil
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string
	Birthday string
	Info     string
	Phone    sql.NullString `gorm:"unique"`

	//处理时区，统一用UTC 0时区的毫秒数
	Ctime int64
	Utime int64
}
