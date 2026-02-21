package interfaces

import (
	"context"
	"time"
)

// LimitRepository adalah kontrak. Middleware tidak peduli ini Redis atau bukan.
type LimitRepository interface {
	Exists(ctx context.Context, key string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Ping(ctx context.Context) error

	// 👇 METODE BARU: Untuk Token Bucket Skala Industri
	TakeToken(ctx context.Context, key string, limit int, burst int, rate float64) (allowed bool, remaining int, err error)
}
