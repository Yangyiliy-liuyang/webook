package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/domain"
	"webook/internal/domain/proctocol"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRouter(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	_, err := a.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: uc.Uid,
		},
	})
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("保存文章数据失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
}
