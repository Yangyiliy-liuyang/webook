package ioc

import (
	"log"
	"os"
	"webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	log.Println(appID, ok)
	if !ok {
		panic("no found in enviroment variable WECHAT_APP_ID")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	log.Println(appID, ok)
	if !ok {
		panic("no found in enviroment variable WECHAT_APP_SECRET")
	}
	return wechat.NewService(appID, appSecret)
}
