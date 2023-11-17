// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	v := ioc.InitGinMiddleware()
	db := ioc.InitDB()
	userDAO := dao.NewGormUserDAO(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	localsmsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, localsmsService)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebService(v, userHandler)
	return engine
}
