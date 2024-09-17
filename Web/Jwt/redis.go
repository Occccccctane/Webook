package Jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

type RedisJWTHandler struct {
	signingMethod jwt.SigningMethod
	client        redis.Cmdable
	rcExpiration  time.Duration
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		signingMethod: jwt.SigningMethodHS256,
		client:        client,
		rcExpiration:  time.Hour * 24 * 7,
	}
}

var JKey = []byte("ppSik8fZfCugefcqWNeh54adKgtN1Fmp")
var RcJwtKey = []byte("ppSik8fZfCugefcqWNeh54adKgtN1FaP")

func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 7天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, rc)
	tokenStr, err := token.SignedString(RcJwtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		//没登录，没有token
		return authCode
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		//token是乱传的
		return ""
	}

	return segs[1]
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) {
	ssid := uuid.New().String()
	err := h.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
		return
	}
	h.SetJWTToken(ctx, uid, ssid)
}

func (h *RedisJWTHandler) SetJWTToken(c *gin.Context, uid int64, ssid string) {

	uc := UserClaims{
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: c.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			//设置15分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(JKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
		return
	}
	c.Header("x-jwt-token", tokenStr)

}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UserClaims)
	return h.client.Set(ctx, fmt.Sprintf("user:ssid:%s", uc.Ssid), "", h.rcExpiration).Err()
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}
