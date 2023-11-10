package service

import "context"

type CodeService struct {
}

// Send 生成一个验证码，发送
func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	return nil
}

// Verify 验证验证码
func (svc *CodeService) Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error) {

}
