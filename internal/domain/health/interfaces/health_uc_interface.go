package interfaces

import (
	"api-guardian/internal/domain/health/dto"
	"context"
)

// HealthUsecase adalah kontrak wajib untuk usecase pengecekan status
type HealthUsecase interface {
	CheckHealth(ctx context.Context) dto.HealthResponse
}
