package ioc

import (
	"log"
	"os"
	"webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("wxa62d4ace23402481")
	log.Println(appID, ok)
	return wechat.NewService(appID)
}
