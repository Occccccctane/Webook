package Failover

import (
	"GinStart/Service/sms"
	"context"
	"sync/atomic"
)

type TimeoutFailOverSMSService struct {
	svcs []sms.Service
	//使用的节点
	idx int32
	//连续超时的计数
	cnt int32
	//切换的阈值，只读，不需要原子操作
	threshold int32
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}
func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	//超过阈值，进行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//	 重置计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplID, args, number...)
	switch err {
	case nil:
		//不超时要重置到0
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		//出现超时，往下+1
		atomic.AddInt32(&t.cnt, 1)
	default:

	}
	return err
}

//	遇到错误但是不是超时，计数可以增加也可以不加
//	如果强调一定是超时，不加
//  如果是EOF类错误，可以考虑直接切换
