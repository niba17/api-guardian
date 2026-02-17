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

var tokenBucketScript = redis.NewScript(`
    local key = KEYS[1]
    local capacity = tonumber(ARGV[1])
    local refill_rate = tonumber(ARGV[2])
    local now = tonumber(ARGV[3])
    
    local bucket = redis.call("HMGET", key, "tokens", "last_refill")
    local last_tokens = tonumber(bucket[1])
    local last_refill = tonumber(bucket[2])

    if last_tokens == nil then
        last_tokens = capacity
        last_refill = now
    end

    local delta = math.max(0, now - last_refill)
    local tokens = math.min(capacity, last_tokens + (delta * refill_rate))

    local allowed = 0
    if tokens >= 1 then
        tokens = tokens - 1
        allowed = 1
    end

    redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
    redis.call("EXPIRE", key, 3600) -- Bucket disimpan 1 jam

    -- Kembalikan angka bulat agar casting di Go aman
    return {allowed, math.floor(tokens)}
`)

func (r *RedisAdapter) TakeToken(ctx context.Context, key string, limit int, burst int, rate float64) (bool, int, error) {
	now := time.Now().Unix()
	res, err := tokenBucketScript.Run(ctx, r.Client, []string{key}, burst, rate, now).Result()
	if err != nil {
		return false, 0, err
	}

	result := res.([]interface{})
	return result[0].(int64) == 1, int(result[1].(int64)), nil
}
