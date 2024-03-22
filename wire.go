//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/jwt"
	"webook/ioc"
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGormInteractiveDAO,
	cache.NewInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//第三方依赖
		ioc.InitLogger, ioc.InitDB, ioc.InitRedis,
		//dao
		dao.NewGormUserDAO, dao.NewGormArticleDAO,
		//cache
		cache.NewRedisUserCache, cache.NewRedisCodeCache, cache.NewArticleRedisCache,
		//repository
		repository.NewCacheUserRepository, repository.NewCodeRepository, repository.NewCachedArticleRepository,
		//service
		ioc.InitSMSService, ioc.InitWechatService,
		service.NewUserService, service.NewCodeService, service.NewArticleService,
		//handler
		jwt.NewRedisJWTHandler,
		web.NewUserHandler, web.NewOAuth2WechatHandler, web.NewArticleHandler,
		ioc.InitGinMiddleware, ioc.InitWebService,
		interactiveSvcSet,
	)
	return gin.Default()
}
