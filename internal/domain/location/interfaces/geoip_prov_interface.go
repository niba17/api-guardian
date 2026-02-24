package interfaces

import (
	"api-guardian/internal/domain/location/dto"
	"context" // 🚀 Tambahkan context
	"io"
)

type GeoIPProvider interface {
	// 🚀 Tambahkan ctx context.Context sebagai parameter pertama
	GetLocation(ctx context.Context, ip string) dto.Location
	io.Closer
}
