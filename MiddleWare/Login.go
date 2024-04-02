package MiddleWare

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginBuilder struct {
}

func (b *LoginBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(context)
		if sess.Get("UserId") == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
