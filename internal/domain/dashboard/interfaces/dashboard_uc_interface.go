package interfaces

import (
	"api-guardian/internal/domain/dashboard/dto" // 👈 Menggunakan DTO
	"context"
)

// DashboardUsecase adalah kontrak untuk logika tampilan dashboard
type DashboardUsecase interface {
	GetStats(ctx context.Context) (dto.DashboardStatsResponse, error)
	GetRecentLogs(ctx context.Context, limit int) ([]dto.RecentLogResponse, error)
}
