package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

// Send 生成一个验证码，发送
func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// Verify 验证验证码
func (svc *CodeService) Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		return false, nil
	}
	return ok, err
}
