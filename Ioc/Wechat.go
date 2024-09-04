package Ioc

import (
	"GinStart/Service/OAuth2/Wechat"
)

func InitWechatService() Wechat.Service {
	// 环境变量获取AppID
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("appID not found")
	//}
	appID := "wx0f0f0f0f0f0f0f0f"
	return Wechat.NewService(appID)
}
