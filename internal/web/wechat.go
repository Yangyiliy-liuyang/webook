package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"webook/internal/domain/proctocol"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	JWTHandler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgC"),
		stateCookieName: "jwt-state",
	}
}

// RegisterRouters 提供两个接口 接口一，构造跳到微信服务的url 接口二，处理跳转回来的请求
func (o *OAuth2WechatHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.OAuth2URL)
	g.Any("/callback", o.Callback)
}

// OAuth2URL 第一次跳转到微信登录
func (o *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	state := uuid.New()
	url, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		resp.SetGeneral(true, 1, "construct url failed")
		resp.SetData(nil)
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		resp.SetGeneral(true, 1, "construct url failed")
	}
	resp.SetGeneral(true, 0, "")
	resp.SetData(url)
	return
}

// Callback 微信回调
// 拉取第三方应用或重定向到第三方应用，带上授权临时票据code
// 通过code加上appid和appsecret换区access_token
// 返回access_code
func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	resp := proctocol.RespGeneral{}
	defer func() {
		ctx.JSON(http.StatusOK, resp)
	}()
	err := o.verifyState(ctx)
	if err != nil {
		resp.SetGeneral(true, 1, "verify state failed")
		resp.SetData(nil)
		return
	}
	code := ctx.Query("code")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		resp.SetGeneral(true, 1, "verify code failed")
		resp.SetData(nil)
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	o.setTokenByJWT(ctx, u.Id)
	resp.SetGeneral(true, 0, "success")
	resp.SetData(nil)
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	token, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return err
	}
	claims := &StateClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return err
	}
	if claims.State != state {
		return fmt.Errorf("state not match %w", err)
	}
	return nil
}

// 放在web层 是认为这是web需要解决的问题，而不是业务
func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	stateClaims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stateClaims)
	tokenString, err := token.SignedString(o.key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}
	//线上环境domain配置成线上环境的域名 测试就测试域名
	// secure 要不要使用https协议 线上环境就配
	ctx.SetCookie(o.stateCookieName, tokenString, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}
