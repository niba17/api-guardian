package app

import (
	"api-guardian/internal/config"
	"api-guardian/internal/database/postgre"
	"api-guardian/internal/database/postgre/migrate"
	"api-guardian/internal/database/postgre/seed"
	"api-guardian/internal/database/redis"
	"api-guardian/internal/delivery/http"
	"api-guardian/internal/delivery/http/handler"
	"api-guardian/internal/delivery/http/middleware"
	authInterfaces "api-guardian/internal/domain/auth/interfaces"
	cacheInterfaces "api-guardian/internal/domain/cache/interfaces"
	dashInterfaces "api-guardian/internal/domain/dashboard/interfaces"
	healthInterfaces "api-guardian/internal/domain/health/interfaces"
	rlInterfaces "api-guardian/internal/domain/rate_limit/interfaces"
	logInterfaces "api-guardian/internal/domain/security_log/interfaces"
	userInterfaces "api-guardian/internal/domain/user/interfaces"
	"api-guardian/internal/infrastructure/auth"
	"api-guardian/internal/proxy"
	"api-guardian/internal/repository"
	"api-guardian/internal/usecase"

	locInterfaces "api-guardian/internal/domain/location/interfaces"

	"api-guardian/internal/database/geoip"
)

func Run(cfg *config.AppConfig) error {
	// 1. Resources & Infrastructure
	dbInstance := postgre.InitDB(cfg.DatabaseDSN)
	rdbInstance := redis.InitRedis(cfg.RedisAddr)
	geoDB := geoip.InitGeoIP(cfg.GeoDBPath)

	if dbInstance != nil && dbInstance.GetDB() != nil {
		migrate.Run(dbInstance.GetDB())
		_ = seed.Run(dbInstance.GetDB())
	}

	lb, err := proxy.NewLoadBalance(cfg)
	if err != nil {
		return err
	}

	// 2. Wiring Repositories (Semua pakai var dan Interface)
	// -----------------------------------------------------
	var rateLimitRepo rlInterfaces.RateLimitRepository = repository.NewRateLimitRepository(rdbInstance)
	var userRepo userInterfaces.UserRepository = repository.NewUserRepository(dbInstance.GetDB())
	var logRepo logInterfaces.SecurityLogRepository = repository.NewSecurityLogRepository(dbInstance.GetDB())
	var cacheRepo cacheInterfaces.CacheRepository = repository.NewCacheRepository(rdbInstance)
	var locRepo locInterfaces.LocationRepository = repository.NewLocationRepository(geoDB)

	// 3. Inisialisasi Usecase (Semua pakai var dan Interface)
	// -----------------------------------------------------
	var authUC authInterfaces.AuthUsecase = usecase.NewAuthUsecase(
		userRepo,
		auth.NewBcryptHasher(),
		auth.NewJWTProvider(cfg.JWTSecret),
	)

	var banUC rlInterfaces.BanUsecase = usecase.NewBanUsecase(
		rateLimitRepo,
		cfg.RateLimit,
	)

	var logUC logInterfaces.SecurityLogUsecase = usecase.NewLogUsecase(
		logRepo,
		banUC,
	)

	var cacheUC cacheInterfaces.CacheUsecase = usecase.NewCacheUsecase(
		cacheRepo,
		cfg.CacheTTL,
	)

	var dashboardUC dashInterfaces.DashboardUsecase = usecase.NewDashboardUsecase(
		cacheRepo,
		logRepo,
	)

	var healthUC healthInterfaces.HealthUsecase = usecase.NewHealthUsecase(
		rateLimitRepo,
	)

	var locUC locInterfaces.LocationUsecase = usecase.NewLocationUsecase(locRepo)

	// 4. Router & Core Handler
	jwtMW := middleware.JWTMiddleware([]byte(cfg.JWTSecret))

	router := http.NewRouter(
		handler.NewAuthHandler(authUC),
		handler.NewDashboardHandler(dashboardUC),
		handler.NewHealthHandler(healthUC),
		jwtMW,
		lb,
	)

	// 5. Wrap Global Middlewares
	finalHandler := wrapGlobalMiddlewares(
		router,
		rateLimitRepo,
		logUC,
		cacheUC,
		banUC,
		locUC,
		cfg,
	)

	// 6. Run Server Lifecycle
	return startServer(finalHandler, cfg.Port, rdbInstance, dbInstance, geoDB)
}
