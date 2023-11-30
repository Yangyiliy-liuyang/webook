package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

var redirectURL = url.PathEscape(`https://meoying.com/oauth2/wechat/callback`)

type service struct {
	appId string
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	//https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}
