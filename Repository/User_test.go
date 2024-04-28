package Repository

import (
	"GinStart/Domain"
	"GinStart/Repository/Cache"
	"GinStart/Repository/Dao"
	cachemock "GinStart/Repository/mocks/cache"
	daomock "GinStart/Repository/mocks/dao"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log"
	"testing"
)

func TestFindByID(t *testing.T) {
	testCase := []struct {
		name         string
		mock         func(ctrl *gomock.Controller) (Dao.UserDao, Cache.UserCache)
		ctx          context.Context
		uid          int64
		ExpectedUser Domain.User
		ExpectedErr  error
	}{
		{
			name: "缓存未命中，查找成功",
			mock: func(ctrl *gomock.Controller) (Dao.UserDao, Cache.UserCache) {
				uid := int64(123)
				d := daomock.NewMockUserDao(ctrl)
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(Domain.User{}, Cache.ErrKeyNotExist)
				d.EXPECT().FindByID(gomock.Any(), uid).Return(Dao.User{
					Id: uid,
					Email: sql.NullString{
						String: "123@qwe.com",
						Valid:  true,
					},
					Password: "asd123",
					Nickname: "aaa",
					Birthday: "1999",
					Info:     "abed",
					Phone: sql.NullString{
						String: "123123123",
						Valid:  true,
					},
				}, nil)
				c.EXPECT().Set(gomock.Any(), Domain.User{
					Id:       uid,
					Email:    "123@qwe.com",
					Password: "asd123",
					Nickname: "aaa",
					Birthday: "1999",
					Info:     "abed",
					Phone:    "123123123",
				}).Return(nil)
				return d, c
			},
			uid: 123,
			ctx: context.Background(),
			ExpectedUser: Domain.User{
				Id:       123,
				Email:    "123@qwe.com",
				Password: "asd123",
				Nickname: "aaa",
				Birthday: "1999",
				Info:     "abed",
				Phone:    "123123123",
			},
			ExpectedErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (Dao.UserDao, Cache.UserCache) {
				uid := int64(123)
				d := daomock.NewMockUserDao(ctrl)
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(Domain.User{}, Cache.ErrKeyNotExist)
				d.EXPECT().FindByID(gomock.Any(), uid).Return(Dao.User{}, ErrUserNotFound)
				return d, c
			},
			uid:          123,
			ctx:          context.Background(),
			ExpectedUser: Domain.User{},
			ExpectedErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			svc := NewCacheUserRepository(ud, uc)

			user, err := svc.FindByID(tc.ctx, tc.uid)
			log.Println(err)
			log.Println(user)

			assert.Equal(t, tc.ExpectedErr, err)
			assert.Equal(t, tc.ExpectedUser, user)
		})
	}
}
