package app

import (
	"api-guardian/internal/config"
	"api-guardian/internal/database/postgre"
	"api-guardian/internal/database/redis"
	"api-guardian/internal/delivery/http"
	"api-guardian/internal/delivery/http/handler"
	"api-guardian/internal/delivery/http/middleware"
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/domain/user"
	"api-guardian/internal/infrastructure/auth"
	"api-guardian/internal/proxy"
	"api-guardian/internal/repository"
	"api-guardian/internal/usecase"

	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

// Run mengoordinasi inisialisasi seluruh komponen aplikasi
func Run(cfg *config.AppConfig) error {

	// 1. Resources & Infrastructure
	rdb := redis.InitRedis(cfg.RedisAddr)
	db := postgre.InitDB(cfg.DatabaseDSN)
	db.AutoMigrate(&user.User{}, &security_log.SecurityLog{})

	// Resource cleanup ditangani di server.go melalui startServer
	geoDB, err := geoip2.Open(cfg.GeoDBPath)
	if err != nil {
		// Kita pakai Warning saja agar aplikasi tidak mati, tapi kita tahu ada yang salah
		log.Warn().
			Str("path", cfg.GeoDBPath).
			Err(err).
			Msg("🌍 GeoIP Database NOT FOUND. Location tracking will be 'Unknown'")
	} else {
		log.Info().Str("path", cfg.GeoDBPath).Msg("🌍 GeoIP Database Loaded Successfully")
	}
	lb, err := proxy.NewLoadBalance(cfg.TargetURLs)
	if err != nil {
		return err // Hentikan jika Load Balancer gagal
	}

	// 2. Wiring Dependencies (Dependency Injection)
	rateLimitRepo := repository.NewRateLimitRepository(rdb)
	userRepo := repository.NewUserRepository(db)
	logRepo := repository.NewSecurityLogRepository(db)
	cacheRepo := repository.NewCacheRepository(rdb)

	// Usecase
	authUC := usecase.NewAuthUsecase(userRepo, auth.NewBcryptHasher(), auth.NewJWTProvider(cfg.JWTSecret))
	banUC := usecase.NewBanUsecase(rateLimitRepo, cfg.RateLimit)
	logUC := usecase.NewLogUsecase(logRepo, banUC) // LogUC butuh BanUC
	cacheUC := usecase.NewCacheUsecase(cacheRepo, cfg.CacheTTL)
	dashboardUC := usecase.NewDashboardUsecase(cacheRepo, logRepo)

	// 3. Router & Core Handler
	jwtMW := middleware.JWTMiddleware([]byte(cfg.JWTSecret))

	// Menggunakan router murni dari delivery/http
	router := http.NewRouter(
		handler.NewAuthHandler(authUC),
		handler.NewDashboardHandler(dashboardUC),
		handler.NewHealthHandler(rateLimitRepo),
		jwtMW,
		lb,
	)

	// 4. Wrap Global Middlewares
	// Memanggil fungsi dari internal/app/middleware.go
	finalHandler := wrapGlobalMiddlewares(router, rateLimitRepo, geoDB, logUC, cacheUC, banUC, cfg)

	// 5. Run Server Lifecycle
	// Memanggil fungsi dari internal/app/server.go (Graceful Shutdown)
	return startServer(finalHandler, cfg.Port, rdb, geoDB)
}
