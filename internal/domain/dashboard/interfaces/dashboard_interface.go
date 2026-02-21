package interfaces

import (
	"api-guardian/internal/domain/security_log"
	"context"
)

// DashboardUsecase adalah kontrak untuk logika tampilan dashboard
type DashboardUsecase interface {
	GetStats(ctx context.Context) (security_log.DashboardStats, error)
	GetRecentLogs(ctx context.Context, limit int) ([]security_log.SecurityLog, error)
}
