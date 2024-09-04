package Wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`

var redirectURL = url.PathEscape(`https://meoying.com/oauth2/wechat/callback`)

type service struct {
	appID string
}

func NewService(appID string) Service {
	return &service{appID: appID}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	state := uuid.New()
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}
