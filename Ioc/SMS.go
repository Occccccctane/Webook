package Ioc

import (
	"GinStart/Service/sms"
	"GinStart/Service/sms/Local"
)

func InitSMSService() sms.Service {
	return Local.NewService()
}
