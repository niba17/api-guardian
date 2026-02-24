package usecase

import (
	cacheIntf "api-guardian/internal/domain/cache/interfaces"
	"context"
	"time"
)

type cacheUsecase struct {
	cacheRepo cacheIntf.CacheRepository
	TTL       time.Duration
}

func NewCacheUsecase(repo cacheIntf.CacheRepository, ttl time.Duration) *cacheUsecase {
	return &cacheUsecase{
		cacheRepo: repo,
		TTL:       ttl,
	}
}

func (u *cacheUsecase) Get(ctx context.Context, key string) (string, error) {
	return u.cacheRepo.Get(ctx, key)
}

// 🚀 FIX: Gunakan 'U' besar pada cacheUsecase dan sesuaikan nama field (cacheRepo & TTL)
func (u *cacheUsecase) Set(ctx context.Context, key string, value interface{}) error {
	return u.cacheRepo.Set(ctx, key, value, u.TTL)
}
