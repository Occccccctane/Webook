package Cache

import (
	"GinStart/Domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *UserCache) Set(ctx context.Context, u Domain.User) error {
	key := c.Key(u.Id)
	// 用JSON进行序列化
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//无论存不存在都重新设置，result字段可以不拿
	return c.cmd.Set(ctx, key, data, c.expiration).Err()

}

func (c *UserCache) Get(ctx context.Context, uid int64) (Domain.User, error) {
	key := c.Key(uid)
	//Redis中数据用JSON来序列化
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return Domain.User{}, err
	}
	var u Domain.User
	err1 := json.Unmarshal([]byte(data), &u)
	return u, err1
}

// Key 格式化成为字符串
func (c *UserCache) Key(uid int64) string {
	//格式没有强制要求,为了和其他的业务隔离开
	//user-info-
	//user.info.
	//user/info/
	//user_info_
	return fmt.Sprintf("user:info:%d", uid)
}
func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd: cmd,
		// 过期时间，专属的可以写死，如果是写通用的缓存机制可以从外部传入
		expiration: time.Minute * 15,
	}
}
