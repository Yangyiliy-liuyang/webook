//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/ioc"
)

var thirdPartySet = wire.NewSet(
	InitDB, InitRedis, InitLog,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//第三方依赖
		thirdPartySet,
		//dao
		dao.NewGormUserDAO, dao.NewGormArticleDAO,
		//cache
		cache.NewRedisUserCache, cache.NewRedisCodeCache,
		//repository
		repository.NewCacheUserRepository, repository.NewCodeRepository, repository.NewCacheArticleRepository,
		//service
		ioc.InitSMSService, InitWechatService,
		service.NewUserService, service.NewCodeService, service.NewArticleService,
		//handler
		ijwt.NewRedisJWTHandler, web.NewUserHandler, web.NewArticleHandler, web.NewOAuth2WechatHandler,
		ioc.InitGinMiddleware, ioc.InitWebService,
	)
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		dao.NewGormArticleDAO,
		repository.NewCacheArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
