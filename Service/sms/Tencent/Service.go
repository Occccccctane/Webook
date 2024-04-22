package Tencent

import (
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// Service   接口sms的具体实现
type Service struct {
	client   *sms.Client
	appID    *string
	signName *string
}

func (s *Service) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appID
	request.SignName = s.signName
	request.TemplateId = &tplID
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(number)
	response, err := s.client.SendSms(request)
	if err != nil {
		fmt.Println("API问题：", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			//不可能进来
			continue
		}
		//取指针
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败,code:%s, msg:%s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toPtrSlice(data []string) []*string {
	ptrSlice := make([]*string, len(data))
	for i, v := range data {
		ptrSlice[i] = &v
	}
	return ptrSlice
}
func NewService(client *sms.Client, appID, signature string) *Service {
	return &Service{
		client:   client,
		appID:    &appID,
		signName: &signature,
	}
}
