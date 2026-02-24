package interfaces

import (
	"api-guardian/internal/domain/location/dto"
	"context"
)

// LocationRepository adalah kontrak untuk mengambil data lokasi dari storage/infra
type LocationRepository interface {
	GetLocationByIP(ctx context.Context, ip string) dto.Location
}
