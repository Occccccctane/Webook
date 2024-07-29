package Failover

import (
	"GinStart/Service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type SMSFailOverService struct {
	svcs []sms.Service

	//	V1的字段
	//当前服务商的下标
	idx uint64
}

func NewSMSFailOverService(svcs []sms.Service) *SMSFailOverService {
	return &SMSFailOverService{
		svcs: svcs,
	}
}

func (f *SMSFailOverService) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplID, args, number...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	//全出问题，很大概率是机房网络出问题，处理不了请求
	return errors.New("轮询所有服务商都发送失败")
}

// 从下标下一位开始轮询，并且出错也轮询
func (f *SMSFailOverService) SendV1(ctx context.Context, tplID string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplID, args, number...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是取消，后者是超时
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询所有服务商都发送失败")
}
