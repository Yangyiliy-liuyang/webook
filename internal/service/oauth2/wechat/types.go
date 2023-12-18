package wechat

import (
	"context"
	"webook/internal/domain"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}
