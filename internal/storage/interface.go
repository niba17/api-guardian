package storage

import (
	"context"
	"time"
)

// LimiterStore adalah kontrak. Middleware tidak peduli ini Redis atau bukan.
type LimiterStore interface {
	// Cek apakah key ada (untuk blacklist)
	Exists(ctx context.Context, key string) (int64, error)

	// Tambah counter (rate limit & violation)
	Incr(ctx context.Context, key string) (int64, error)

	// Set expire waktu
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)

	// Set key dengan value (untuk banned)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Ping untuk health check
	Ping(ctx context.Context) error
}
