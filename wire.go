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
	"webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//第三方依赖
		ioc.InitDB, ioc.InitRedis,
		//dao
		dao.NewUserDAO,
		//cache
		cache.NewUserCache, cache.NewCodeCache,
		//repository
		repository.NewUserRepository, repository.NewCodeRepository,
		//service
		ioc.InitSMSService, service.NewUserService, service.NewCodeService,
		//handler
		web.NewUserHandler, ioc.InitGinMiddleware, ioc.InitWebService,
	)
	return gin.Default()
}
