//go:build !k8s

package Config

// Config 本地启动
var Config = config{
	DB: DBConfig{
		DSN: "root:aaa@tcp(localhost:3306)/ginstart",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
