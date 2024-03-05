package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"webook/internal/domain"
)

var redirectURL = url.PathEscape(`https://meoying.com/oauth2/wechat/callback`)

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appId string) Service {
	return &service{
		appId:  appId,
		client: http.DefaultClient,
		//appSecret: appSecret,
	}
}

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type Result struct {
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unianid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// 通过code发起调用获取accessToken
func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", s.appId, s.appSecret, code)
	Req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(Req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var req Result
	err = json.NewDecoder(httpResp.Body).Decode(req)
	if err != nil {
		//json 失败
		return domain.WechatInfo{}, err
	}
	if req.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 %s", err)
	}
	return domain.WechatInfo{}, err
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	//https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}
