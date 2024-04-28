package Ioc

import (
	"GinStart/Config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{Addr: Config.Config.Redis.Addr})
}
