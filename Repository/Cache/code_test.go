package Cache

import (
	"GinStart/Repository/mocks/cache/redismock"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	keyfunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCase := []struct {
		name        string
		mock        func(ctrl *gomock.Controller) redis.Cmdable
		ctx         context.Context
		biz         string
		phone       string
		code        string
		ExpectedErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmd.SetErr(nil)

				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyfunc("test", "1333333333")}, []any{"200"}).Return(cmd)
				return c
			},
			ctx:         context.Background(),
			biz:         "test",
			phone:       "1333333333",
			code:        "200",
			ExpectedErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis错误"))

				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyfunc("test", "1333333333")}, []any{"200"}).Return(cmd)
				return c
			},
			ctx:         context.Background(),
			biz:         "test",
			phone:       "1333333333",
			code:        "200",
			ExpectedErr: errors.New("redis错误"),
		},
		{
			name: "验证码存在，没有过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-2))
				cmd.SetErr(errors.New("验证码存在，但是没有过期时间"))

				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyfunc("test", "1333333333")}, []any{"200"}).Return(cmd)
				return c
			},
			ctx:         context.Background(),
			biz:         "test",
			phone:       "1333333333",
			code:        "200",
			ExpectedErr: errors.New("验证码存在，但是没有过期时间"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))
				cmd.SetErr(ErrCodeSendToMany)

				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyfunc("test", "1333333333")}, []any{"200"}).Return(cmd)
				return c
			},
			ctx:         context.Background(),
			biz:         "test",
			phone:       "1333333333",
			code:        "200",
			ExpectedErr: ErrCodeSendToMany,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cc := tc.mock(ctrl)
			c := NewCodeCache(cc)

			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.ExpectedErr, err)
		})
	}
}
