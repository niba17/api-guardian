package usecase

import (
	"api-guardian/internal/domain/location/dto"
	"api-guardian/internal/domain/location/interfaces"
	"context"
)

type locationUsecase struct {
	locRepo interfaces.LocationRepository // 🚀 GANTI: Sekarang bergantung pada Repo
}

// 🚀 GANTI: Constructor menerima parameter Repository
func NewLocationUsecase(repo interfaces.LocationRepository) interfaces.LocationUsecase {
	return &locationUsecase{
		locRepo: repo,
	}
}

func (u *locationUsecase) GetLocationByIP(ctx context.Context, ip string) dto.Location {
	// Usecase hanya fokus ke bisnis logik, urusan ambil data serahkan ke Repo
	return u.locRepo.GetLocationByIP(ctx, ip)
}
