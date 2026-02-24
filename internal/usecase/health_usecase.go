package usecase

import (
	"api-guardian/internal/domain/health/dto" // 👈 Tambahkan import DTO baru kita
	healthInterfaces "api-guardian/internal/domain/health/interfaces"
	rateLimitInterfaces "api-guardian/internal/domain/rate_limit/interfaces"
	"context"
	"time"
)

type healthUsecase struct {
	redisRepo rateLimitInterfaces.RateLimitRepository
}

func NewHealthUsecase(redisRepo rateLimitInterfaces.RateLimitRepository) healthInterfaces.HealthUsecase {
	return &healthUsecase{
		redisRepo: redisRepo,
	}
}

// 🚀 Ubah return type menjadi dto.HealthResponse
func (u *healthUsecase) CheckHealth(ctx context.Context) dto.HealthResponse {
	redisStatus := "Connected"

	if u.redisRepo != nil {
		if err := u.redisRepo.Ping(ctx); err != nil {
			redisStatus = "Disconnected: " + err.Error()
		}
	} else {
		redisStatus = "Not Configured"
	}

	// 🚀 Return menggunakan dto.HealthResponse
	return dto.HealthResponse{
		System:          "Healthy",
		RedisConnection: redisStatus,
		CircuitBreaker:  "Closed (Normal)",
		Time:            time.Now().UTC().Format(time.RFC3339),
	}
}
