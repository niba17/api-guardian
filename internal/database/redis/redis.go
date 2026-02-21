package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient = redis.Client

func InitRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            addr,
		MaxRetries:      10,
		MinRetryBackoff: 1 * time.Second,
		DialTimeout:     15 * time.Second,
		PoolSize:        50,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Fail to connect Redis, system run without cache")
	} else {
		log.Info().Str("addr", addr).Msg("Redis Infrastructure Ready")
	}

	return rdb
}
