package localsms

import (
	"context"
	"log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Send(ctx context.Context, tplId string, args []string, phone ...string) error {
	log.Println(" 验证码: ", args)
	return nil
}
