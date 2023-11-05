package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

// CheckLogin 登录校验
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// todo 注册下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			//中断
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		now := time.Now()
		const UpdateTimeKey = "update_time"
		val := sess.Get(UpdateTimeKey)
		LastUpdateTime, ok := val.(time.Time)
		if val != nil || !ok || now.Sub(LastUpdateTime) > time.Minute {
			//第一次进来
			sess.Set(UpdateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				//打印日志
				println("err")
			}
		}
	}
}
