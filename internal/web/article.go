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
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		ID int64 `json:"id"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := a.svc.Withdraw(ctx, req.ID, uc.Uid)
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("撤回文章数据失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "")
	resp.SetData(nil)
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		ID      int64 `json:"id"`
		Title   string
		Content string
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	artId, err := a.svc.Publish(ctx, domain.Article{
		Id:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("发布文章数据失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(artId)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		ID      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	artId, err := a.svc.Save(ctx, domain.Article{
		Id:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("保存文章数据失败", logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(artId)
}
