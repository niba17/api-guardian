package main

import (
	"api-guardian/internal/config"
	"api-guardian/internal/database" // 👈 Import database
	"api-guardian/internal/handler"
	"api-guardian/internal/middleware"
	"api-guardian/internal/proxy"
	"api-guardian/internal/usecase"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	setupLogger()

	// 1. Load Config
	cfg := config.Load()

	// 2. Setup Dependencies
	// a. Redis
	cacheUC := usecase.NewCacheUsecase(cfg.RedisClient, cfg.CacheTTL)

	// b. Postgres Database (🔥 NEW)
	// Pastikan variable DatabaseDSN sudah ada di config.go Bos yang terakhir
	db := database.InitDB(cfg.DatabaseDSN)

	// c. GeoIP Database
	geoDB, err := geoip2.Open(cfg.GeoDBPath)
	if err != nil {
		log.Warn().Err(err).Str("path", cfg.GeoDBPath).Msg("Fail to load GeoIP Database")
	} else {
		defer geoDB.Close()
		log.Info().Str("path", cfg.GeoDBPath).Msg("GeoIP Database loaded")
	}

	// 3. Setup Load Balancer
	lb, err := proxy.NewLoadBalancer(cfg.TargetURLs)
	if err != nil {
		log.Fatal().Err(err).Msg("Target URL Error")
	}

	// 4. Setup Middleware Chain
	var secureHandler http.Handler = lb

	secureHandler = middleware.SmartCache(cacheUC, secureHandler)
	secureHandler = middleware.APIKeyValidator(cfg.APIKeys, secureHandler)
	secureHandler = middleware.BasicWAF(secureHandler)
	secureHandler = middleware.PrometheusMiddleware(secureHandler)
	secureHandler = middleware.RateLimiter(cfg.Storage, cfg.WhitelistIPs, secureHandler)

	// 👇 Pass 'db' ke AuditLogger
	secureHandler = middleware.AuditLogger(geoDB, cfg.RedisClient, db, secureHandler)

	// 5. Setup Router
	mux := http.NewServeMux()

	mux.HandleFunc("/status", handler.HealthCheck(cfg.Storage))
	mux.Handle("/metrics", promhttp.Handler())

	// --- ENDPOINT DASHBOARD ---
	// 👇 Pass 'db' ke Dashboard Handler
	dashHandler := handler.NewDashboardHandler(cfg.RedisClient, db)

	mux.HandleFunc("/api/dashboard/stats", dashHandler.GetDashboardStats)
	mux.HandleFunc("/api/dashboard/logs", dashHandler.GetRecentLogs)

	// Endpoint Utama
	mux.Handle("/", secureHandler)

	// 6. Server Run
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Info().
			Str("port", cfg.Port).
			Int("backends", len(cfg.TargetURLs)).
			Msg("API Guardian is standing guard (Powered by PostgreSQL)")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info().Msg("Server shutdown complete")
}

func setupLogger() {
	// Logger file ini sekarang HANYA untuk debug aplikasi/system error
	// Audit log keamanan sudah masuk ke PostgreSQL
	fileLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // Ganti nama jadi app.log biar gak kecampur
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)
	log.Logger = log.Output(multi)
}
