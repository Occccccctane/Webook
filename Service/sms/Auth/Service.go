package Auth

import (
	"GinStart/Service/sms"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceSMS struct {
	svc sms.Service
	key []byte
}

func (s *AuthServiceSMS) Send(ctx context.Context, tplToken string, args []string, number ...string) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, claims.Tpl, args, number...)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
}
