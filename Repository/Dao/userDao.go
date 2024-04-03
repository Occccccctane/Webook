package Dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	EmailUniqueErr    = errors.New("邮箱唯一错误")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
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

func (dao *UserDao) EmailSearch(context context.Context, email string) (User, error) {
	var user User
	err1 := dao.db.WithContext(context).Where("email=?", email).First(&user).Error
	if err1 != nil {
		return user, err1
	}
	return user, nil
}

func (dao *UserDao) Update(user User, password string) (User, error) {
	//更新信息
	var newUser User
	newUser = user
	newUser.Password = password
	now := time.Now().UnixMilli()
	newUser.Ctime = now

	err1 := dao.db.Save(&newUser).Error
	if err1 != nil {
		return user, err1
	}
	return newUser, nil
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	//处理时区，统一用UTC 0时区的毫秒数
	Ctime int64
	Utime int64
}
