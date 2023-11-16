package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func InitWebService(funcs []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(funcs...)
	userHdl.RegisterRouter(server)
	return server
}
func InitGinMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3030"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//允许前端访问后端响应中带的头部
		ExposeHeaders:    []string{"X-Jwt-Token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//if strings.HasPrefix(origin,"http://localhost") {
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "公司域名.com")
		},
		MaxAge: 12 * time.Hour,
	}),
		// todo 限流
		(&middleware.LoginJWTMilddlewareBuilder{}).CheckLoginJWT(),
	}
}
