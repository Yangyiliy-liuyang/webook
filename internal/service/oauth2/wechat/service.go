package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"webook/internal/domain"
	"webook/pkg/logger"
)

var redirectURL = url.PathEscape(`https://meoying.com/oauth2/wechat/callback`)

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	l         logger.Logger
}

func NewService(appId string, appSecret string, l logger.Logger) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		l:         l,
	}
}

type WechatResult struct {
	AccessToken  string `json:"access_token"`  //接口调用凭证
	ExpiresIn    int64  `json:"expires_in"`    // 接口调用凭证的超时时间，单位：秒
	RefreshToken string `json:"refresh_token"` // 刷新access_token的凭证
	OpenId       string `json:"openid"`        // 用户的id
	Scope        string `json:"scope"`         // 用户的权限
	UnionId      string `json:"unianid"`       // 用户的id 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段
	//错误信息
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

// 通过code发起调用获取accessToken
func (s *service) VerifyCode(ctx *gin.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", s.appId, s.appSecret, code)
	Req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(Req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var req WechatResult
	err = json.NewDecoder(httpResp.Body).Decode(req)
	if err != nil {
		return domain.WechatInfo{}, fmt.Errorf("解析json失败 %s", err)
	}
	if req.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 %s", err)
	}
	return domain.WechatInfo{
		OpenId:  req.OpenId,
		UnionId: req.UnionId,
	}, nil
}

func (s *service) AuthURL(ctx *gin.Context, state string) (string, error) {
	// 微信接口链接 https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}
