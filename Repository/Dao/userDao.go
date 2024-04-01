package Dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var EmailUniqueErr = errors.New("邮箱唯一错误")

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

	if mysl, ok := err.(*mysql.MySQLError); ok {
		const uniqueErrNum uint16 = 1062
		if mysl.Number == uniqueErrNum {
			return EmailUniqueErr
		}
	}
	return err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	//处理时区，统一用UTC 0时区的毫秒数
	Ctime int64
	Utime int64
}
