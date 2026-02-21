package usecase

import (
	cacheIntf "api-guardian/internal/domain/cache/interfaces"
	dashboardIntf "api-guardian/internal/domain/dashboard/interfaces"
	"api-guardian/internal/domain/security_log"
	logIntf "api-guardian/internal/domain/security_log/interfaces"
	"context"
)

type dashboardUsecase struct {
	cacheRepo cacheIntf.CacheRepository
	logRepo   logIntf.SecurityLogRepository
}

// Constructor sekarang meminta cacheRepo, bukan lagi rdb mentah
func NewDashboardUsecase(cacheRepo cacheIntf.CacheRepository, repo logIntf.SecurityLogRepository) dashboardIntf.DashboardUsecase {
	return &dashboardUsecase{
		cacheRepo: cacheRepo,
		logRepo:   repo,
	}
}

func (u *dashboardUsecase) GetStats(ctx context.Context) (security_log.DashboardStats, error) {
	// 💡 Nanti logika cek Redis bisa pakai: u.cacheRepo.Get(ctx, "dashboard:stats")
	return u.logRepo.GetStats()
}

func (u *dashboardUsecase) GetRecentLogs(ctx context.Context, limit int) ([]security_log.SecurityLog, error) {
	return u.logRepo.GetRecent(limit)
}
