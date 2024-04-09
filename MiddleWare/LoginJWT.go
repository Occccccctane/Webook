package MiddleWare

import (
	Handler "GinStart/Web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJWTBuilder struct {
}

func (b *LoginJWTBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		//约定token在Authorization的Bearer一起请求
		authCode := context.GetHeader("Authorization")
		if authCode == "" {
			//没登录，没有token
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			//token是乱传的
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		var uc Handler.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, uc, func(token *jwt.Token) (interface{}, error) {
			return Handler.JWTKey, nil
		})

		//token不对是伪造的 || token没解析出来 || token是非法的或是过期的
		if err != nil || token == nil || !token.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		//如果剩余时间少于50秒，刷新
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 15))
			tokenStr, err := token.SignedString(Handler.JWTKey)
			context.Header("x-jwt-token", tokenStr)
			if err != nil {
				fmt.Println(err)
			}
		}
		context.Set("user", uc) //将其放入上下文中
	}
}
