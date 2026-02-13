package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisAdapter membungkus library asli redis agar sesuai interface kita
type RedisAdapter struct {
	Client *redis.Client
}

// NewRedisAdapter adalah constructor
func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{Client: client}
}

func (r *RedisAdapter) Exists(ctx context.Context, key string) (int64, error) {
	return r.Client.Exists(ctx, key).Result()
}

func (r *RedisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

func (r *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(ctx, key, expiration).Result()
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisAdapter) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
