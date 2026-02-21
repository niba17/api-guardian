package usecase

import (
	"api-guardian/internal/domain/rate_limit/interfaces" // 👈 1. Arahkan ke folder interfaces
	"context"
	"time"
)

type BanUsecase struct {
	limitRepo     interfaces.LimitRepository
	maxViolations int
}

// 👈 3. Sesuaikan parameter constructor
func NewBanUsecase(repo interfaces.LimitRepository, maxViolations int) *BanUsecase {
	return &BanUsecase{
		limitRepo:     repo,
		maxViolations: maxViolations,
	}
}

// ExecuteAutoBan mengevaluasi apakah IP harus diblokir
func (u *BanUsecase) ExecuteAutoBan(ctx context.Context, ip string) {
	violationKey := "violation:" + ip
	banKey := "blacklist:" + ip

	// 1. Catat/Tambah pelanggaran
	count, _ := u.limitRepo.Incr(ctx, violationKey)

	// 2. Set expiry 1 jam untuk pelanggaran pertama (agar reset otomatis)
	if count == 1 {
		_, _ = u.limitRepo.Expire(ctx, violationKey, 1*time.Hour)
	}

	// 3. 🚀 UBAH INI: Gunakan variabel dinamis dari config, bukan angka 5 lagi!
	if count >= int64(u.maxViolations) {
		_ = u.limitRepo.Set(ctx, banKey, "BANNED", 24*time.Hour)
	}
}
