package Failover

import (
	smsmock "GinStart/Service/mocks/sms"
	"GinStart/Service/sms"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSMSFailOverService_Send(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(ctrl *gomock.Controller) []sms.Service
		ExpectedErr error
	}{
		{
			name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
		},
		{
			//轮询
			name: "第二次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc1 := smsmock.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
		},
		{
			//轮询失败
			name: "全部轮询失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc1 := smsmock.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return []sms.Service{svc0, svc1}
			},
			ExpectedErr: errors.New("轮询所有服务商都发送失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewSMSFailOverService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "123", []string{"1234"}, "1234")
			assert.Equal(t, tc.ExpectedErr, err)
		})
	}
}
