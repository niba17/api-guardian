package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient = redis.Client

// 🚀 1. BUAT JEMBATAN LOGGER: Agar log internal go-redis masuk ke Zerolog
type redisLogger struct{}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	// Tangkap semua log internal Redis dan ubah jadi format Warning di Zerolog
	log.Warn().Msgf("go-redis internal: "+format, v...)
}

func InitRedis(addr string) *redis.Client {
	// 🚀 2. PASANG JEMBATANNYA: Bungkam logger bawaan!
	redis.SetLogger(&redisLogger{})

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		MaxRetries:   1,               // Batas eksekusi perintah
		DialTimeout:  2 * time.Second, // Batas nunggu socket
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     50,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Fail to connect Redis, system run without cache/rate-limit")
	} else {
		log.Info().Str("addr", addr).Msg("Redis Infrastructure Ready")
	}

	return rdb
}
