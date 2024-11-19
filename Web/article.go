package Handler

import (
	"GinStart/Domain"
	"GinStart/Service"
	"GinStart/Web/Jwt"
	"GinStart/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleHandler struct {
	svc Service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc Service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}
func (h *ArticleHandler) RegisterRoute(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(Jwt.UserClaims)
	id, err := h.svc.Save(ctx, Domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: Domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		h.l.Error("保存文章失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "保存成功",
		Data: id,
	})
}
