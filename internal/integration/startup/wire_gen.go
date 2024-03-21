// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

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

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := InitLog()
	v := ioc.InitGinMiddleware(handler, logger)
	db := InitDB()
	userDAO := dao.NewGormUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	localsmsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, localsmsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := InitWechatService(logger)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDAO := dao.NewGormArticleDAO(db)
	articleRepository := repository.NewCachedArticleRepository(articleDAO)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, logger)
	engine := ioc.InitWebService(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

func InitArticleHandler(articleDAO dao.ArticleDAO) *web.ArticleHandler {
	articleRepository := repository.NewCachedArticleRepository(articleDAO)
	articleService := service.NewArticleService(articleRepository)
	logger := InitLog()
	articleHandler := web.NewArticleHandler(articleService, logger)
	return articleHandler
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitDB, InitRedis, InitLog,
)
