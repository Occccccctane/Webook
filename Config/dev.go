//go:build !k8s

package Config

// Config 本地启动
var Config = config{
	DB: DBConfig{
		DSN: "root:aaa@tcp(localhost:13316)/Gin",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
