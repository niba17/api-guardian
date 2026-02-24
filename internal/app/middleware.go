package app

import (
	"api-guardian/internal/config"
	"api-guardian/internal/delivery/http/middleware"

	// 🚀 1. Import interface-nya saja
	cacheInterfaces "api-guardian/internal/domain/cache/interfaces"
	locInterfaces "api-guardian/internal/domain/location/interfaces"
	rateLimitInterfaces "api-guardian/internal/domain/rate_limit/interfaces"
	logInterfaces "api-guardian/internal/domain/security_log/interfaces"
	"net/http"
	// ❌ HAPUS import "github.com/oschwald/geoip2-golang"
)

func wrapGlobalMiddlewares(
	router http.Handler,
	RatelimitRepo rateLimitInterfaces.RateLimitRepository,
	// ❌ HAPUS geo *geoip2.Reader dari sini
	logUC logInterfaces.SecurityLogUsecase,
	cacheUC cacheInterfaces.CacheUsecase,
	banUC rateLimitInterfaces.BanUsecase,
	locUC locInterfaces.LocationUsecase,
	cfg *config.AppConfig,
) http.Handler {

	h := router

	// 1. LAPISAN PALING DALAM: Cache melayani data HANYA JIKA tamu sudah lolos sekuriti
	h = middleware.Cache(cacheUC, h)
	h = middleware.APIKeyValidator(cfg.APIKeys)(h)
	h = middleware.BasicWAF(cfg.WhitelistIPs)(h)
	h = middleware.RateLimit(RatelimitRepo, banUC, cfg.WhitelistIPs, cfg.RefillRate, cfg.BurstCapacity, h)

	// 5. OBSERVABILITY: Logger wajib mencatat semua request masuk
	h = middleware.Logger(locUC, logUC, h)

	h = middleware.CORSMiddleware(cfg.AllowedOrigins)(h)

	// 7. LAPISAN PALING LUAR: Prometheus (Menghitung waktu dari ujung ke ujung)
	return middleware.PrometheusMiddleware(h)
}
