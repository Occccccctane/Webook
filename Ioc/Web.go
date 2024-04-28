package Ioc

import (
	"GinStart/MiddleWare"
	Handler "GinStart/Web"
	"GinStart/pkg/middleware/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(middleware []gin.HandlerFunc, userHdl *Handler.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middleware...)
	userHdl.RegisterRoute(server)
	return server
}

func InitMiddleWare(client redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//跨域
		(&MiddleWare.CrossDomain{}).CrossDomainHandler(),
		//限流
		ratelimit.NewBuilder(client, time.Second, 100).Build(),
		//登录校验
		(&MiddleWare.LoginJWTBuilder{}).CheckLogin(),
	}
}
