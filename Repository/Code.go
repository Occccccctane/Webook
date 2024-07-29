package Repository

import (
	"GinStart/Repository/Cache"
	"context"
)

var (
	ErrVerifySendToMany = Cache.ErrVerifySendToMany
	ErrCodeSendToMany   = Cache.ErrCodeSendToMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CacheCodeRepository struct {
	cache Cache.CodeCache
}

func NewCodeRepository(c Cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (c *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}
func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
