package usecase

import (
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/domain/security_log/interfaces"
	"context"
)

type LogUsecase struct {
	logRepo interfaces.SecurityLogRepository
	banUC   *BanUsecase
}

func NewLogUsecase(l interfaces.SecurityLogRepository, b *BanUsecase) *LogUsecase {
	return &LogUsecase{logRepo: l, banUC: b}
}

func (u *LogUsecase) LogAndEvaluate(ctx context.Context, logData *security_log.SecurityLog, reason string) {
	// 1. Catat aktivitas
	_ = u.logRepo.Save(logData)

	// 2. Jika ada indikasi serangan, delegasikan ke BanUsecase
	if reason != "" && reason != "IP Permanently Banned" {
		u.banUC.ExecuteAutoBan(ctx, logData.IP)
	}
}
