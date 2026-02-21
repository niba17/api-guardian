package app

import (
	"api-guardian/internal/config"
	"api-guardian/internal/delivery/http/middleware"
	"api-guardian/internal/domain/rate_limit/interfaces"
	"api-guardian/internal/usecase"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

func wrapGlobalMiddlewares(
	router http.Handler,
	RatelimitRepo interfaces.LimitRepository,
	geo *geoip2.Reader,
	logUC *usecase.LogUsecase,
	cacheUC *usecase.CacheUsecase,
	banUC *usecase.BanUsecase,
	cfg *config.AppConfig,
) http.Handler {

	h := router

	// 1. LAPISAN PALING DALAM: Cache melayani data HANYA JIKA tamu sudah lolos sekuriti
	h = middleware.Cache(cacheUC, h)

	// 2. LAPISAN SEKURITI 3: Cek Kunci (API Key)
	if len(cfg.APIKeys) > 0 {
		h = middleware.APIKeyValidator(cfg.APIKeys)(h)
	}

	// 3. LAPISAN SEKURITI 2: WAF (Cek Payload Jahat)
	h = middleware.BasicWAF(h)

	// 4. LAPISAN SEKURITI 1: Rate Limit & Ban (Garda Depan Satpam)
	h = middleware.RateLimit(RatelimitRepo, banUC, cfg.WhitelistIPs, cfg.RefillRate, cfg.BurstCapacity, h)

	// 5. OBSERVABILITY: Logger wajib mencatat semua request masuk sebelum difilter
	h = middleware.Logger(geo, logUC, h)

	// 6. STANDAR WEB: CORS
	h = middleware.CORSMiddleware(cfg.AllowedOrigins)(h)

	// 7. LAPISAN PALING LUAR: Prometheus (Menghitung waktu dari ujung ke ujung)
	return middleware.PrometheusMiddleware(h)
}
