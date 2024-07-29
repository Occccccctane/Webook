package Cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCodeCache_Set_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCase := []struct {
		name        string
		before      func(t *testing.T)
		after       func(t *testing.T)
		ctx         context.Context
		biz         string
		phone       string
		code        string
		ExpectedErr error
	}{
		{
			name: "设置成功",
			before: func(t *testing.T) {
				//	不需要
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:13333333333"

				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)

				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)

			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "13333333333",
			code:        "123456",
			ExpectedErr: nil,
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13333333333"
				err := rdb.Set(ctx, key, "654321", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:13333333333"

				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)

				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)

			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "13333333333",
			code:        "123456",
			ExpectedErr: ErrCodeSendToMany,
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13333333333"
				err := rdb.Set(ctx, key, "654321", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:13333333333"

				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)

			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "13333333333",
			code:        "123456",
			ExpectedErr: errors.New("验证码存在，但是没有过期时间"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			c := NewCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.ExpectedErr, err)
		})
	}
}
