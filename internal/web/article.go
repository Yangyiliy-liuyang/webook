package web

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
	"webook/internal/domain"
	"webook/internal/domain/proctocol"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

type ArticleHandler struct {
	svc     service.ArticleService
	intrSvc service.InteractiveService
	l       logger.Logger
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger, intrSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		intrSvc: intrSvc,
		biz:     "article",
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

	// 读者接口
	pub := g.Group("/pub")
	pub.GET("/detail:id", a.PubDetail)
	pub.GET("/like", a.Like)
	pub.POST("/collection", a.Collection)
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

func (a *ArticleHandler) PubDetail(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer ctx.JSON(http.StatusOK, resp)
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

		ReadCnt    int64 `json:"read_cnt"`
		LikeCnt    int64 `json:"like_cnt"`
		CollectCnt int64 `json:"collect_cnt"`
		Liked      bool  `json:"liked"`
		Collected  bool  `json:"collected"`
	}
	var data article
	str := ctx.Param("id")
	artId, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	var (
		eg   errgroup.Group
		art  domain.Article
		intr domain.Interactive
	)
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	eg.Go(func() error {
		var er error
		art, er = a.svc.GetPubByArtId(ctx, artId, uc.Uid)
		return er
	})

	eg.Go(func() error {
		var er error
		intr, er = a.intrSvc.GetIntrByArtId(ctx, a.biz, artId, uc.Uid)
		return er
	})
	if err := eg.Wait(); err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("获取文章详情数据失败", logger.Int64("uid", art.Author.Id), logger.Int64("id", art.Id), logger.Error(err))
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err = a.intrSvc.IncrReadCnt(ctx, a.biz, art.Id)
		if err != nil {
			a.l.Error("文章阅读量增加失败", logger.Int64("id", art.Id), logger.Error(err))
		}
	}()
	data = article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status.ToUint8(),
		// 需要
		AuthorName: art.Author.Name,
		Ctime:      art.Ctime,
		Utime:      art.Utime,

		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(data)
}

func (a *ArticleHandler) Like(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		ArtId int64 `json:"art_id"`
		Like  bool  `json:"like"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	var err error
	if req.Like {
		err = a.intrSvc.Like(ctx, a.biz, req.ArtId, uc.Uid)
	} else {
		err = a.intrSvc.CancelLike(ctx, a.biz, req.ArtId, uc.Uid)
	}
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("点赞失败", logger.Int64("uid", uc.Uid), logger.Int64("id", req.ArtId), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(nil)

}

func (a *ArticleHandler) Collection(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	type Req struct {
		ArtId int64 `json:"art_id"`
		Cid   int64 `json:"cid"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.SetGeneral(true, http.StatusBadRequest, "参数错误")
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	err := a.intrSvc.AddCollectionItem(ctx, a.biz, req.ArtId, req.Cid, uc.Uid)
	if err != nil {
		resp.SetGeneral(true, http.StatusInternalServerError, "系统内部错误")
		a.l.Error("收藏失败", logger.Int64("uid", uc.Uid), logger.Int64("id", req.ArtId), logger.Error(err))
		return
	}
	resp.SetGeneral(true, http.StatusOK, "ok")
	resp.SetData(nil)

}
