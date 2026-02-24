package interfaces

import (
	"api-guardian/internal/domain/security_log"
	"context"
)

type SecurityLogUsecase interface {
	// 🚀 Gunakan nama ini agar sinkron dengan implementasi Usecase Bos
	LogAndEvaluate(ctx context.Context, logData *security_log.SecurityLog, reason string)
}
