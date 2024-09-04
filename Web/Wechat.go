package Handler

import (
	"GinStart/Service/OAuth2/Wechat"
	"github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
	svc Wechat.Service
}

func NewOAuth2WechatHandler(svc Wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (o *OAuth2WechatHandler) RegisterRout(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {

	val, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(200, Result{
			Code: 500,
			Msg:  "构造跳转URL失败 ",
		})
		return
	}
	ctx.JSON(200, Result{
		Code: 200,
		Data: val,
		Msg:  "构造跳转URL成功 ",
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}
