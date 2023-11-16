package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
)

func InitSMSService() *localsms.Service {
	return localsms.NewService()
	//InitTencentSMSService
}

func InitTencentSMSService() sms.Service {
	//todo tencentSMS
	return nil
}
