package usecase

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheUsecase menangani logika bisnis caching
type CacheUsecase struct {
	Redis *redis.Client
	TTL   time.Duration
}

func NewCacheUsecase(rdb *redis.Client, ttl time.Duration) *CacheUsecase {
	return &CacheUsecase{
		Redis: rdb,
		TTL:   ttl,
	}
}

// Get mengambil data cache berdasarkan key
func (u *CacheUsecase) Get(ctx context.Context, key string) (string, error) {
	return u.Redis.Get(ctx, key).Result()
}

// Set menyimpan data ke cache dengan durasi TTL
func (u *CacheUsecase) Set(ctx context.Context, key string, value []byte) error {
	return u.Redis.Set(ctx, key, value, u.TTL).Err()
}
