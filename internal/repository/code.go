package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(codeCache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: codeCache,
	}
}
func (repo *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}