package Local

import (
	"context"
	"log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	log.Println("验证码是", args)
	return nil
}
