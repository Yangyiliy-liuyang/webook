package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	// 创作者接口
	g.POST("/list", a.List)
	g.GET("/detail:id", a.Detail)
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

func (a *ArticleHandler) List(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		Limit  int `json:"limit"`
		Offset int `bson:"offset"`
	}
	type article struct {
		Id         int64  `json:"id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Status     uint8  `json:"status"`
		AuthorId   int64  `json:"author_id"`
		AuthorName string `json:"author_name"`
		Ctime      int64  `json:"ctime"`
		Utime      int64  `json:"utime"`
	}
	var data []article
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	arts, err := a.svc.GetByAuthor(ctx, req.Limit, req.Offset, uc.Uid)
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("获取文章列表数据失败", logger.Int("offset", req.Offset),
			logger.Int("limit", req.Limit),
			logger.Int64("uid", uc.Uid), logger.Error(err))
		return
	}
	data = slice.Map[domain.Article, article](arts, func(idx int, src domain.Article) article {
		return article{
			Id:      src.Id,
			Title:   src.Title,
			Content: src.Content,
			Status:  src.Status.ToUint8(),
			// 不需要Author作者信息
			//Ctime: src.Ctime,
			//Utime: src.Utime,
			Ctime: src.Ctime,
			Utime: src.Utime,
		}
	})
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(data)
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type article struct {
		Id         int64  `json:"id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Status     uint8  `json:"status"`
		AuthorId   int64  `json:"author_id"`
		AuthorName string `json:"author_name"`
		Ctime      int64  `json:"ctime"`
		Utime      int64  `json:"utime"`
	}
	var data article
	str := ctx.Param("id")
	artId, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	art, err := a.svc.GetByArtId(ctx, artId)
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("获取文章详情数据失败", logger.Int64("uid", art.Author.Id), logger.Int64("id", art.Id), logger.Error(err))
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	if art.Author.Id != uc.Uid {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("没有权限，获取文章详情数据失败", logger.Int64("uid", art.Author.Id), logger.Int64("id", art.Id), logger.Error(err))
		return
	}
	data = article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status.ToUint8(),
		Ctime:   art.Ctime,
		Utime:   art.Utime,
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(data)
}
