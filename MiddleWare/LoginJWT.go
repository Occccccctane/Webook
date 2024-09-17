package MiddleWare

import (
	ijwt "GinStart/Web/Jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTBuilder struct {
	ijwt.Handler
}

func NewLoginJWTBuilder(hdl ijwt.Handler) *LoginJWTBuilder {
	return &LoginJWTBuilder{
		Handler: hdl,
	}
}
func (b *LoginJWTBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			return
		}

		//约定token在Authorization的Bearer一起请求
		tokenStr := b.ExtractToken(context)

		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JKey, nil
		})

		//token不对是伪造的 || token没解析出来 || token是非法的或是过期的
		if err != nil || token == nil || !token.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//校验完token后再访问redis可降低一些无效的访问场景
		err = b.CheckSession(context, uc.Ssid)
		if err != nil {
			return
		}
		//严格做法
		if err != nil {
			// token无效或是redis出问题
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//可以兼容redis出现问题，redis出现问题可以继续访问
		//需要做好监控有没有 error
		//if cnt > 0 {
		// token无效
		//	context.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		context.Set("user", uc) //将其放入上下文中
	}
}
