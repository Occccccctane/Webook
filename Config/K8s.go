//go:build k8s

package Config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(ginstart-mysql:3308)/Gin",
	},
	Redis: RedisConfig{
		Addr: "ginstart-redis:6379",
	},
}
