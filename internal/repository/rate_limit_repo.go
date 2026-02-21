package repository

import (
	"api-guardian/internal/domain/rate_limit/interfaces"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimitRepo adalah implementasi tertutup (private struct)
type RateLimitRepo struct {
	client *redis.Client
}

// NewRateLimitRepository merakit Redis dan mengembalikan kontrak interface
func NewRateLimitRepository(client *redis.Client) interfaces.LimitRepository {
	return &RateLimitRepo{client: client}
}

func (r *RateLimitRepo) Exists(ctx context.Context, key string) (int64, error) {
	return r.client.Exists(ctx, key).Result()
}

func (r *RateLimitRepo) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RateLimitRepo) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

func (r *RateLimitRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RateLimitRepo) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// --- Token Bucket Script Skala Industri ---
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

	return {allowed, math.floor(tokens)}
`)

func (r *RateLimitRepo) TakeToken(ctx context.Context, key string, rate_limit int, burst int, rate float64) (bool, int, error) {
	now := time.Now().Unix()
	res, err := tokenBucketScript.Run(ctx, r.client, []string{key}, burst, rate, now).Result()
	if err != nil {
		return false, 0, err
	}

	result := res.([]interface{})
	return result[0].(int64) == 1, int(result[1].(int64)), nil
}
