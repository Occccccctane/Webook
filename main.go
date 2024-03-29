package main

import (
	"GinStart/Handler"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {
	hdl := Handler.NewUserHandler()

	c := gin.Default()
	c.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,允许所有域名，比较危险
		//AllowedOrigins: []string{"http://localhost:3000","http://aaa"},枚举允许的域名

		AllowCredentials: true,                     //是否允许带cookie等用户凭据，正常都需要允许
		AllowedHeaders:   []string{"content-type"}, // 允许的请求头

		//AllowedMethods: []string{"POST"}, 允许的请求方法，最好不配置

		//允许字符串的检查方法，如果传入的字符串包含相关的字段则放行
		AllowOriginFunc: func(origin string) bool {

			if strings.Contains(origin, "localhost") { //判断包含该字段
				// if strings.HasPrefix(origin, "http://localhost")判定包含前缀
				return true
			}
			return strings.Contains(origin, "xxx.com") //返回一个表达式，上面用于判断是否是本机调试，下面用于判断是否是公司的域名
		},

		MaxAge: 12 * time.Hour, //检测时间长度
	}))

	hdl.RegisterRoute(c)
	c.Run(":8080")
}
