package repository

import (
	"api-guardian/internal/storage"
	"context"
	"time"
)

type BlacklistRepository interface {
	IncrViolation(ctx context.Context, ip string) (int64, error)
	SetBan(ctx context.Context, ip string, duration time.Duration) error
	IsBanned(ctx context.Context, ip string) (bool, error)
	SetExpire(ctx context.Context, key string, duration time.Duration) error
}

type redisBlacklistRepo struct {
	storage storage.LimiterStore
}

func NewBlacklistRepository(s storage.LimiterStore) BlacklistRepository {
	return &redisBlacklistRepo{storage: s}
}

func (r *redisBlacklistRepo) IncrViolation(ctx context.Context, ip string) (int64, error) {
	return r.storage.Incr(ctx, "violation:"+ip)
}

func (r *redisBlacklistRepo) SetBan(ctx context.Context, ip string, duration time.Duration) error {
	return r.storage.Set(ctx, "blacklist:"+ip, "banned", duration)
}

func (r *redisBlacklistRepo) IsBanned(ctx context.Context, ip string) (bool, error) {
	val, err := r.storage.Exists(ctx, "blacklist:"+ip)
	return val > 0, err
}

func (r *redisBlacklistRepo) SetExpire(ctx context.Context, key string, duration time.Duration) error {
	_, err := r.storage.Expire(ctx, key, duration)
	return err
}
