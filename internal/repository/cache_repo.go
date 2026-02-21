package repository

import (
	"api-guardian/internal/domain/cache/interfaces"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCacheRepo adalah struct tertutup (private)
type redisCacheRepo struct {
	client *redis.Client
}

// NewCacheRepository merakit Redis dan mengembalikan interface CacheRepository
func NewCacheRepository(client *redis.Client) interfaces.CacheRepository {
	return &redisCacheRepo{
		client: client,
	}
}

func (r *redisCacheRepo) Get(ctx context.Context, key string) (string, error) {
	// Mengambil data dari Redis berdasarkan key
	return r.client.Get(ctx, key).Result()
}

func (r *redisCacheRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Menyimpan data ke Redis dengan waktu kedaluwarsa (TTL)
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisCacheRepo) Delete(ctx context.Context, key string) error {
	// Menghapus data dari Redis
	return r.client.Del(ctx, key).Err()
}
