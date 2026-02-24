package usecase

import (
	rlInterfaces "api-guardian/internal/domain/rate_limit/interfaces"
	"api-guardian/internal/domain/security_log"
	logInterfaces "api-guardian/internal/domain/security_log/interfaces"
	"context"
)

// 🚀 1. GANTI: Jadi huruf kecil (private).
// Orang luar tidak boleh bikin struct ini secara manual, harus lewat NewLogUsecase.
type securityLogUsecase struct {
	logRepo logInterfaces.SecurityLogRepository
	banUC   rlInterfaces.BanUsecase
}

// 🚀 2. GANTI: Return type-nya adalah Interface dari Domain, bukan struct pointer.
func NewLogUsecase(l logInterfaces.SecurityLogRepository, b rlInterfaces.BanUsecase) logInterfaces.SecurityLogUsecase {
	return &securityLogUsecase{
		logRepo: l,
		banUC:   b,
	}
}

// 🚀 3. SESUAIKAN: Receiver-nya pakai nama struct yang huruf kecil tadi
func (u *securityLogUsecase) LogAndEvaluate(ctx context.Context, logData *security_log.SecurityLog, reason string) {
	// 1. Catat aktivitas
	_ = u.logRepo.Save(logData)

	// 2. Jika ada indikasi serangan, delegasikan ke BanUsecase
	if reason != "" && reason != "IP Blacklisted" {
		u.banUC.ExecuteAutoBan(ctx, logData.IP)
	}
}
