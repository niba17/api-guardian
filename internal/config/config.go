package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Port         string
	TargetURLs   []string
	APIKeys      []string
	WhitelistIPs []string
	CacheTTL     time.Duration
	GeoDBPath    string
	DatabaseDSN  string
	RedisAddr    string
	JWTSecret    string
}

func Load() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Info().Msg(".env tidak ditemukan, menggunakan Environment Variables sistem")
	}

	return &AppConfig{
		Port:         getEnv("PORT", "8080"),
		TargetURLs:   parseCSV(getEnv("TARGET_URL", "http://localhost:8081")),
		APIKeys:      parseCSV(getEnv("API_KEYS", "")),
		WhitelistIPs: parseCSV(getEnv("WHITELIST_IPS", "")),
		CacheTTL:     parseDuration(getEnv("CACHE_TTL", "60s")),
		GeoDBPath:    getEnv("GEOIP_DB_PATH", "configs/geoip/GeoLite2-City.mmdb"),
		DatabaseDSN:  getEnv("DATABASE_DSN", ""),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:    getEnv("JWT_SECRET", "rahasia-negara-bos-jangan-disebar"),
	}
}

// Helper sederhana agar Load() tetap bersih
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	list := strings.Split(s, ",")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}
	return list
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
		return 60 * time.Second
	}
	return v
}
