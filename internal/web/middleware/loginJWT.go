package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	ijwt "webook/internal/web/jwt"
)

type LoginJWTMilddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMilddlewareBuilder(hdl ijwt.Handler) *LoginJWTMilddlewareBuilder {
	return &LoginJWTMilddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMilddlewareBuilder) CheckLoginJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			return
		}
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			//token不对伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// 解析出来了 但是是非法的过期的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//expireTime := uc.ExpiresAt
		//if expireTime.Before(time.Now()) {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//if uc.UserAgent != ctx.GetHeader("User-Agent") {
		//	// todo 埋点
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		////过期时间小于10分钟刷新
		//if expireTime.Sub(time.Now()) < time.Minute*10 {
		//	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString(web.JWTKey)
		//	ctx.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		//过期时间没有刷新，已登录
		//		log.Println(err)
		//	}
		//}
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			//做好Redis崩溃的预警
			return
		}
		ctx.Set("user", uc)
	}
}
