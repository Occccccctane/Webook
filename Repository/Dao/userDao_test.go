package Dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDiver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(T *testing.T) *sql.DB
		ctx         context.Context
		user        User
		ExpectedErr error
	}{
		{
			name: "插入成功",
			mock: func(T *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(T, err)
				//创建结果集
				mockres := sqlmock.NewResult(123, 1)
				//要求传入一个sql的正则表达式,返回结果集
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockres)
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "aaa",
			},
			ExpectedErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(T *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(T, err)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysqlDiver.MySQLError{Number: 1062}) //返回mysql唯一索引冲突错误码1062
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "aaa",
			},
			ExpectedErr: EmailUniqueErr,
		},
		{
			//数据库错误
			name: "插入失败",
			mock: func(T *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(T, err)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("数据库错误")) //返回mysql唯一索引冲突错误码1062
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "aaa",
			},
			ExpectedErr: errors.New("数据库错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			ud := NewUserDao(db)

			err1 := ud.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.ExpectedErr, err1)

		})
	}
}
