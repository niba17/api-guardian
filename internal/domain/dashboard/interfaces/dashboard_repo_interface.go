package interfaces

import (
	"api-guardian/internal/domain/dashboard" // 👈 Import domain dashboard yang baru
	"api-guardian/internal/domain/security_log"
	"context"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (dashboard.Dashboard, error)

	GetRecentLogs(ctx context.Context, limit int) ([]security_log.SecurityLog, error)
}
