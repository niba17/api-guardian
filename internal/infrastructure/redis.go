package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func InitRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            addr,
		MaxRetries:      10,
		MinRetryBackoff: 1 * time.Second,
		DialTimeout:     15 * time.Second,
		PoolSize:        50,
	})

	// Pastikan koneksi benar-benar siap
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Gagal koneksi ke Redis, sistem berjalan tanpa cache")
	} else {
		log.Info().Str("addr", addr).Msg("Redis Infrastructure Ready")
	}

	return rdb
}
