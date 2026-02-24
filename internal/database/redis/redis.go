package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisInstance interface {
	redis.Cmdable
	Close() error
}

// 🛑 1. BUNGKAM LOGGER BAWAAN: Jangan biarkan go-redis cerewet per request!
type silentLogger struct{}

func (l *silentLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	// Dikosongkan dengan sengaja. Kita tidak mau log bawaannya.
}

func InitRedis(addr string) *redis.Client {
	// Pasang penyumbat logger
	redis.SetLogger(&silentLogger{})

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		MaxRetries:   1,               // Jangan coba ulang berkali-kali kalau memang mati
		DialTimeout:  1 * time.Second, // Cepat menyerah agar tidak membebani CPU
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     50,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Fail to connect Redis on startup. System is in Fail-Closed mode!")
	} else {
		log.Info().Str("addr", addr).Msg("Redis Infrastructure Ready")
	}

	// 🚀 2. NYALAKAN SMART MONITOR DI BACKGROUND
	go StartRedisMonitor(rdb)

	return rdb
}

// 🧠 START SMART MONITOR (Background Ping)
func StartRedisMonitor(rdb *redis.Client) {
	isOnline := true // Anggap awalnya hidup

	for {
		time.Sleep(3 * time.Second) // Cek setiap 3 detik

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err := rdb.Ping(ctx).Err()
		cancel()

		if err != nil {
			// Kalau sebelumnya online, lalu tiba-tiba mati -> Teriak SEKALI SAJA!
			if isOnline {
				log.Error().Msg("🔴 ALERT: REDIS CONNECTION LOST! WAF entering Fail-Closed mode. Standing by for reconnect...")
				isOnline = false
			}
		} else {
			// Kalau sebelumnya mati, lalu tiba-tiba hidup -> Lapor SEKALI SAJA!
			if !isOnline {
				log.Info().Msg("🟢 RECOVERY: REDIS IS BACK ONLINE! WAF fully operational.")
				isOnline = true
			}
		}
	}
}
