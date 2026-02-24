package usecase

import (
	cacheIntf "api-guardian/internal/domain/cache/interfaces"
	"api-guardian/internal/domain/dashboard/dto"
	dashboardIntf "api-guardian/internal/domain/dashboard/interfaces"
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

// --- 1. GET STATS (Mapper Entity -> DTO) ---
func (uc *dashboardUsecase) GetStats(ctx context.Context) (dto.DashboardStatsResponse, error) {
	// 1. Ambil data mentah dari Repository
	rawStats, err := uc.logRepo.GetStats()
	if err != nil {
		return dto.DashboardStatsResponse{}, err
	}

	// 2. Bungkus ke dalam DTO dengan presisi 100%
	return dto.DashboardStatsResponse{
		TotalRequests: rawStats.TotalRequests,
		TotalBlocked:  rawStats.TotalBlocked,
		TotalSuccess:  rawStats.TotalSuccess,
		UniqueIPs:     rawStats.UniqueIPs,
		AvgLatency:    rawStats.AvgLatency,
		TopIPs:        rawStats.TopIPs,
	}, nil
}

// --- 2. GET RECENT LOGS (Mapper Array Entity -> Array DTO) ---
func (uc *dashboardUsecase) GetRecentLogs(ctx context.Context, limit int) ([]dto.RecentLogResponse, error) {
	// 1. Ambil array data mentah dari Repository
	rawLogs, err := uc.logRepo.GetRecent(limit)
	if err != nil {
		return nil, err
	}

	// 2. Siapkan keranjang kosong untuk DTO
	var result []dto.RecentLogResponse

	// 3. Konversi satu per satu dari Entity ke DTO
	for _, log := range rawLogs {
		result = append(result, dto.RecentLogResponse{
			ID:            log.ID,
			Timestamp:     log.Timestamp,
			IP:            log.IP,
			Country:       log.Country,
			City:          log.City, // 🚀 Transfer City
			Method:        log.Method,
			Path:          log.Path,
			Status:        log.Status,
			IsBlocked:     log.IsBlocked,
			Latency:       log.Latency,
			ThreatType:    log.ThreatType,    // 🚀 Transfer ThreatType
			ThreatDetails: log.ThreatDetails, // 🚀 TRANSFER REASON KE SINI!
			Browser:       log.Browser,       // 🚀 Transfer Browser
			OS:            log.OS,            // 🚀 Transfer OS
			IsBot:         log.IsBot,         // 🚀 Transfer IsBot
		})
	}

	return result, nil
}
