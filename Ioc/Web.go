package Ioc

import (
	"GinStart/MiddleWare"
	Handler "GinStart/Web"
	"GinStart/pkg/limiter"
	"GinStart/pkg/middleware/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(middleware []gin.HandlerFunc, userHdl *Handler.UserHandler, wechatHdl *Handler.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middleware...)
	userHdl.RegisterRoute(server)
	wechatHdl.RegisterRout(server)
	return server
}

func InitMiddleWare(client redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//跨域
		(&MiddleWare.CrossDomain{}).CrossDomainHandler(),
		//限流
		ratelimit.NewBuilder(limiter.NewRedisSlideWindowsLimiter(client, time.Second, 1000)).Build(),
		//登录校验
		(&MiddleWare.LoginJWTBuilder{}).CheckLogin(),
	}
}
