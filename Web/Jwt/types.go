package Jwt

import "github.com/gin-gonic/gin"

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetLoginToken(ctx *gin.Context, uid int64)
	SetJWTToken(c *gin.Context, uid int64, ssid string)
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
}
