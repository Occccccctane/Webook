package Ratelimit

import (
	"GinStart/Service/sms"
	"GinStart/pkg/limiter"
	"context"
	"errors"
)

var errRateLimit = errors.New("触发限流")

type RateLimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errRateLimit
	}
	return r.svc.Send(ctx, tplID, args, number...)
}

func NewRateLimitSMSService(svc sms.Service, l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "limiter",
	}
}
