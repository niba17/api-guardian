package interfaces

import (
	"api-guardian/internal/domain/security_log"
	"context"
)

// DashboardRepository adalah pintu gerbang data untuk Dashboard
type DashboardRepository interface {
	GetStats(ctx context.Context) (security_log.DashboardStats, error)
	GetRecentLogs(ctx context.Context, limit int) ([]security_log.SecurityLog, error)
}
