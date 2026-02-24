package interfaces

import (
	"api-guardian/internal/domain/dashboard"
	"api-guardian/internal/domain/security_log"
)

type SecurityLogRepository interface {
	Save(log *security_log.SecurityLog) error
	GetStats() (dashboard.Dashboard, error)
	GetRecent(limit int) ([]security_log.SecurityLog, error)
}
