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
		panic("[WECHAT_APP_ID] err ")
	}
	return wechat.NewService(appID)
}
