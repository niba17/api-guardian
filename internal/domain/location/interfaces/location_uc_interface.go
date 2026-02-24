package interfaces

import (
	"api-guardian/internal/domain/location/dto"
	"context"
)

// LocationUsecase adalah kontrak untuk logika bisnis terkait pencarian lokasi
type LocationUsecase interface {
	GetLocationByIP(ctx context.Context, ip string) dto.Location
}
