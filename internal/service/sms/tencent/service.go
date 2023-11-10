package tencent

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type Service struct {
	client   *sms.Client
	appId    *string
	SignName *string
}

func NewService(client *sms.Client, appId string, SignName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		SignName: &SignName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.SignName
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(numbers)
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code != nil || *(status.Code) != "ok" {
			//发送失败
			return fmt.Errorf("短信发送失败，code:%s,message:%s", *status.Code, *status.Message)
		}
	}
	return nil
}
