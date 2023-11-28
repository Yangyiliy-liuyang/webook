package web

import "github.com/gin-gonic/gin"

type OAuth2WechatHandler struct {
}

func (o *OAuth2WechatHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.OAuth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {

}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}
