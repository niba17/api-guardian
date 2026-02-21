package interfaces

import (
	"context"
	"time"
)

// CacheRepository adalah kontrak untuk penyimpanan sementara (Caching)
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}
