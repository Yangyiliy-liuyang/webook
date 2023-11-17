package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: codeCache,
	}
}
func (repo *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
