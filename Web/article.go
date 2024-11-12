package Handler

import "github.com/gin-gonic/gin"

type ArticleHandler struct {
}

func (h *ArticleHandler) RegisterRoute(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	return
}
