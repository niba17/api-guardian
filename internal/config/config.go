package config

import (
	"api-guardian/internal/storage"
	"context"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Port         string
	TargetURLs   []string
	APIKeys      []string
	WhitelistIPs []string
	CacheTTL     time.Duration
	GeoDBPath    string // 🟢 BARU: Supaya lokasi DB dinamis
	Storage      storage.LimiterStore
	RedisClient  *redis.Client
}

func Load() *AppConfig {
	// 1. Coba load .env (Tapi jangan panic kalau gagal, karena di Docker kita pakai Env Vars langsung)
	if err := godotenv.Load(); err != nil {
		log.Info().Msg(".env not found, using System Environment Variables")
	}

	// 2. Setup Redis
	rdb := connectRedisClient()

	// 3. Parsing Target URL (Comma Separated)
	rawTargets := getEnv("TARGET_URL", "http://localhost:8081")
	targetList := strings.Split(rawTargets, ",")
	for i := range targetList {
		targetList[i] = strings.TrimSpace(targetList[i])
	}

	return &AppConfig{
		Port:         getEnv("PORT", "8080"),
		TargetURLs:   targetList,
		APIKeys:      strings.Split(getEnv("API_KEYS", ""), ","),
		WhitelistIPs: strings.Split(getEnv("WHITELIST_IPS", ""), ","),
		CacheTTL:     parseDuration(getEnv("CACHE_TTL", "60s")),

		// 🟢 BARU: Default path mengarah ke folder configs/geoip yang baru kita buat
		GeoDBPath: getEnv("GEOIP_DB_PATH", "configs/geoip/GeoLite2-City.mmdb"),

		Storage:     storage.NewRedisAdapter(rdb),
		RedisClient: rdb,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseDuration(d string) time.Duration {
	v, err := time.ParseDuration(d)
	if err != nil {
		log.Warn().Msgf("Wrong duration format: %s. using default 60s", d)
		return 60 * time.Second
	}
	return v
}

func connectRedisClient() *redis.Client {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	rdb := redis.NewClient(&redis.Options{
		Addr:            addr,
		MaxRetries:      10,               // Naikkan jumlah percobaan otomatis
		MinRetryBackoff: 1 * time.Second,  // Jeda antar percobaan lebih santai
		DialTimeout:     15 * time.Second, // Kasih waktu lebih lama buat ngetok pintu
		PoolSize:        50,
		ConnMaxIdleTime: 5 * time.Minute,
	})

	// --- PROTEKSI STARTUP: Lebih sabar menunggu ---
	var err error
	for i := 0; i < 7; i++ { // Coba 7 kali
		// Naikkan dari 5s ke 10s 👇
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Info().Msg("Redis connected")
			return rdb
		}
		log.Warn().Msgf("Waiting Redis (retry %d/7)... error: %v", i+1, err)
		time.Sleep(3 * time.Second) // Tunggu 3 detik sebelum coba lagi
	}

	log.Error().Err(err).Msg("Fail to connect Redis (System running without cache)")
	return rdb
}
