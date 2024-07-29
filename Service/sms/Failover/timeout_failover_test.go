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

func TestTimeoutFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(ctrl *gomock.Controller) []sms.Service
		threshold   int32
		idx         int32
		cnt         int32
		ExpectedErr error
		ExpectedIdx int32
		ExpectedCnt int32
	}{
		{
			name: "没有触发切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			idx:         0,
			cnt:         12,
			threshold:   15,
			ExpectedErr: nil,
			ExpectedIdx: 0,
			ExpectedCnt: 0,
		},
		{
			name: "触发切换，成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc1 := smsmock.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:         0,
			cnt:         15,
			threshold:   15,
			ExpectedErr: nil,
			ExpectedIdx: 1,
			ExpectedCnt: 0,
		},
		{
			name: "触发切换，失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc1 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return []sms.Service{svc0, svc1}
			},
			idx:         1,
			cnt:         15,
			threshold:   15,
			ExpectedIdx: 0,
			ExpectedErr: errors.New("发送失败"),
		},
		// todo 修改mock
		{
			name: "触发切换，超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmock.NewMockService(ctrl)
				svc1 := smsmock.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				//svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			idx:         1,
			cnt:         15,
			threshold:   15,
			ExpectedIdx: 0,
			ExpectedCnt: 1,
			ExpectedErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailOverSMSService(tc.mock(ctrl), tc.threshold)
			svc.cnt = tc.cnt
			svc.idx = tc.idx

			err := svc.Send(context.Background(), "1234", []string{"12", "34"}, "12341234")
			assert.Equal(t, tc.ExpectedErr, err)
			assert.Equal(t, tc.ExpectedCnt, svc.cnt)
			assert.Equal(t, tc.ExpectedIdx, svc.idx)
		})
	}
}
