package usecase

import (
	"api-guardian/internal/domain/rate_limit/interfaces"
	"context"
	"time"
)

type banUsecase struct {
	limitRepo     interfaces.RateLimitRepository
	maxViolations int
}

func NewBanUsecase(repo interfaces.RateLimitRepository, maxViolations int) *banUsecase {
	return &banUsecase{
		limitRepo:     repo,
		maxViolations: maxViolations,
	}
}

// 🚀 GANTI: Namanya sekarang ExecuteAutoBan (Exported) untuk memenuhi Interface
func (u *banUsecase) ExecuteAutoBan(ctx context.Context, ip string) error {
	violationKey := "violation:" + ip
	banKey := "blacklist:" + ip

	count, _ := u.limitRepo.Incr(ctx, violationKey)

	if count == 1 {
		_, _ = u.limitRepo.Expire(ctx, violationKey, 1*time.Hour)
	}

	if count >= int64(u.maxViolations) {
		return u.limitRepo.Set(ctx, banKey, "BANNED", 24*time.Hour)
	}
	return nil
}
