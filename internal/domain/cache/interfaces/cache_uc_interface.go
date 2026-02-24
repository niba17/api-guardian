package interfaces

import "context"

type CacheUsecase interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
}
