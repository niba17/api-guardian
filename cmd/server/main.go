package main

import (
	"api-guardian/internal/config"
	"api-guardian/internal/database"
	"api-guardian/internal/domain"

	httpDelivery "api-guardian/internal/delivery/http"
	"api-guardian/internal/delivery/http/handler"
	"api-guardian/internal/delivery/http/middleware" // Middleware HTTP

	"api-guardian/internal/infrastructure"
	"api-guardian/internal/proxy"
	"api-guardian/internal/repository"
	"api-guardian/internal/storage"
	"api-guardian/internal/usecase"
	"api-guardian/pkg/logger"

	"net/http"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.Setup()
	cfg := config.Load()

	// 1. Infrastructure
	rdb := infrastructure.InitRedis(cfg.RedisAddr)
	db := database.InitDB(cfg.DatabaseDSN)
	db.AutoMigrate(&domain.User{}, &domain.SecurityLog{})

	limiterStore := storage.NewRedisAdapter(rdb)
	geoDB, _ := geoip2.Open(cfg.GeoDBPath)
	lb, _ := proxy.NewLoadBalancer(cfg.TargetURLs)

	// 2. Repository
	logRepo := repository.NewSecurityLogRepository(db)
	blRepo := repository.NewBlacklistRepository(limiterStore)
	userRepo := repository.NewUserRepository(db)

	// 3. Usecase
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	cacheUC := usecase.NewCacheUsecase(rdb, cfg.CacheTTL)

	// 4. Handler
	authHandler := handler.NewAuthHandler(authUC)
	dashHandler := handler.NewDashboardHandler(rdb, logRepo)
	healthHandler := handler.NewHealthHandler(limiterStore)

	// 5. Router (Mux)
	authMW := middleware.AuthMiddleware([]byte(cfg.JWTSecret))

	// secureChain di sini adalah Load Balancer (Proxy) yang akan menangani 404
	// atau request yang tidak di-handle oleh router API.
	var fallbackHandler http.Handler = lb

	mux := httpDelivery.NewRouter(
		authHandler,
		dashHandler,
		healthHandler,
		authMW,
		fallbackHandler,
	)

	// 6. GLOBAL SECURITY MIDDLEWARE CHAIN (Pengganti ApplySecurityMiddleware)
	// Kita bungkus Router (mux) dengan layer keamanan.
	// Urutan Eksekusi: Request Masuk -> WAF -> RateLimit -> Router

	var finalHandler http.Handler = mux

	// Layer 3: WAF (Paling Dalam - Cek isi request)
	finalHandler = middleware.BasicWAF(finalHandler)

	// Layer 2: Rate Limiter & Blacklist (Tengah - Cek IP/Daftar Cekal)
	// 👈 SEKARANG INI DI LUAR WAF
	finalHandler = middleware.RateLimiter(limiterStore, cfg.WhitelistIPs, finalHandler)

	// Layer 1: Audit Logger (Paling Luar - Rekam Segalanya)
	finalHandler = middleware.AuditLogger(geoDB, logRepo, blRepo, finalHandler)

	// Layer 0: Sisanya (Cache, CORS, dll)
	finalHandler = middleware.SmartCache(cacheUC, finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)
	finalHandler = middleware.PrometheusMiddleware(finalHandler)

	// 7. Server Run
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Info().Msgf("🛡️ API Guardian Standing Guard on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
