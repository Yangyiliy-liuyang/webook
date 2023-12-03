package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{svc: svc}
}

// RegisterRouters 提供两个接口 接口一，构造跳到微信服务的url 接口二，处理跳转回来的请求
func (o *OAuth2WechatHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.OAuth2URL)
	g.Any("/callback", o.Callback)
}

// OAuth2URL 第一次跳转到微信登录
func (o *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	url, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转微信登录Url失败",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "构造URL成功",
		Data: url,
	})
}

// Callback 微信回调
// 拉取第三方应用或重定向到第三方应用，带上授权临时票据code
// 通过code加上appid和appsecret换区access_token
// 返回access_code
func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}
