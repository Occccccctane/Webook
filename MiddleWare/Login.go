package MiddleWare

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginBuilder struct {
}

func (b *LoginBuilder) CheckLogin() gin.HandlerFunc {
	//为成功将时间类型转成字节流，注册这个类型
	gob.Register(time.Now())
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(context)
		userId := sess.Get("UserId")
		if userId == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//设置更新时间
		now := time.Now()
		const updateTimeKey = "update_time"

		uptime := sess.Get(updateTimeKey)
		lastUpdateTime, ok := uptime.(time.Time)

		//判断条件
		//1.刷新时间不为nil，即第一次登录
		//2.上次刷新时间，断言不成功，如果刷新时间为nil也会断言不成功
		//3.现在时间减上次时间大于设置的刷新时间
		if uptime == nil || !ok || now.Sub(lastUpdateTime) > time.Second*30 {
			sess.Set("update_time", now)
			sess.Set("UserId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
