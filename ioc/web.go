package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/internal/web/middleware"
	"webook/pkg/logger"
)

func InitWebService(funcs []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(funcs...)
	userHdl.RegisterRouter(server)
	wechatHdl.RegisterRouters(server)
	return server
}

func InitGinMiddleware(hdl ijwt.Handler, l *logger.ZapLogger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins: []string{"http://localhost:3030"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			//允许前端访问后端响应中带的头部
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				//if strings.HasPrefix(origin,"http://localhost") {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "公司域名.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		// todo 限流
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AcessLog) {
			l.Debug("", logger.Field{
				Key:   "req",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMilddlewareBuilder(hdl).CheckLoginJWT(),
	}
}
