package Handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JwtHandler struct {
	signingMethod jwt.SigningMethod
	refreshKey    []byte
}

func NewJwtHandler() *JwtHandler {
	return &JwtHandler{
		signingMethod: jwt.SigningMethodHS256,
		refreshKey:    []byte("ppSik8fZfCugefcqWNeh54adKgtN1FmP"),
	}
}

var JWTKey = []byte("ppSik8fZfCugefcqWNeh54adKgtN1Fmp")

func (h *JwtHandler) SetJWTToken(c *gin.Context, uid int64) {
	err := h.SetRefreshToken(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
		return
	}
	uc := UserClaims{
		Uid:       uid,
		UserAgent: c.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			//设置15分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
		return
	}
	c.Header("x-jwt-token", tokenStr)

}

func (h *JwtHandler) SetRefreshToken(ctx *gin.Context, uid int64) error {
	rc := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 7天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, rc)
	tokenStr, err := token.SignedString(h.refreshKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
