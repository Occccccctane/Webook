package sms

import "context"

// Service 发送短信的抽象，屏蔽不同供应商的差别
type Service interface {
	Send(ctx context.Context, tplID string,
		args []string, number ...string) error
}
