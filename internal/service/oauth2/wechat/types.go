package wechat

import (
	"github.com/gin-gonic/gin"
	"webook/internal/domain"
)

type Service interface {
	AuthURL(ctx *gin.Context, state string) (string, error)
	VerifyCode(ctx *gin.Context, code string) (domain.WechatInfo, error)
}
