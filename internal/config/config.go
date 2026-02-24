package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Port           string
	TargetURLs     []string
	APIKeys        []string
	WhitelistIPs   []string
	AllowedOrigins []string
	CacheTTL       time.Duration
	GeoDBPath      string
	DatabaseDSN    string
	RedisAddr      string
	JWTSecret      string
	RateLimit      int
	RefillRate     float64
	BurstCapacity  int
	CBMaxRequests  uint32
	CBIntervalSec  time.Duration
	CBTimeoutSec   time.Duration
	CBMinRequests  uint32
	CBFailRatio    float64
}

func Load() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Info().Msg(".env not found, using Environment Variables system")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "api_guardian")
	dbSSL := getEnv("DB_SSLMODE", "disable")
	dbTZ := getEnv("DB_TIMEZONE", "UTC")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSL, dbTZ)

	return &AppConfig{
		Port:           getEnv("PORT", "8080"),
		TargetURLs:     parseCSV(getEnv("TARGET_URL", "https://www.google.com")),
		APIKeys:        parseCSV(getEnv("API_KEYS", "")),
		WhitelistIPs:   parseCSV(getEnv("WHITELIST_IPS", "")),
		AllowedOrigins: parseCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:5173")),
		CacheTTL:       parseDuration(getEnv("CACHE_TTL", "2s")),
		GeoDBPath:      getEnv("GEO_DB_PATH", "configs/geoip/GeoLite2-City.mmdb"),
		DatabaseDSN:    dsn,
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", "rahasia-negara-bos-jangan-disebar"),
		RateLimit:      parseInt(getEnv("RATE_LIMIT", "5")),
		RefillRate:     parseFloat(getEnv("REFILL_RATE", "0.5")),
		BurstCapacity:  parseInt(getEnv("BURST_CAPACITY", "10")),
		CBMaxRequests:  parseUint32(getEnv("CB_MAX_REQUESTS", "1")),
		CBIntervalSec:  time.Duration(parseInt(getEnv("CB_INTERVAL_SEC", "10"))),
		CBTimeoutSec:   time.Duration(parseInt(getEnv("CB_TIMEOUT_SEC", "30"))),
		CBMinRequests:  parseUint32(getEnv("CB_MIN_REQUESTS", "3")),
		CBFailRatio:    parseFloat(getEnv("CB_FAIL_RATIO", "0.6")),
	}
}

// --- HELPERS ---

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

func parseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return v
}

// 🚀 Helper Baru untuk uint32 (Circuit Breaker)
func parseUint32(s string) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}
