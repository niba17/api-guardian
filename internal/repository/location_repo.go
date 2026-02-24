package repository

import (
	"api-guardian/internal/domain/location/dto"
	"api-guardian/internal/domain/location/interfaces"
	"context"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// 🚀 STRUCT KECIL: Privat untuk package repository
type locationRepository struct {
	db *geoip2.Reader
}

// 🚀 CONSTRUCTOR: Menerima *geoip2.Reader (beton) dari app.go, mengembalikan Interface
func NewLocationRepository(db *geoip2.Reader) interfaces.LocationRepository {
	return &locationRepository{
		db: db,
	}
}

// 🚀 METHOD: Mengandung logika pencarian MaxMind
func (r *locationRepository) GetLocationByIP(ctx context.Context, ip string) dto.Location {
	loc := dto.Location{Country: "Unknown", City: "Unknown"}

	if r.db == nil || ip == "" {
		return loc
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return loc
	}

	record, err := r.db.City(parsedIP)
	if err == nil {
		if c, ok := record.Country.Names["en"]; ok {
			loc.Country = c
		}
		if c, ok := record.City.Names["en"]; ok {
			loc.City = c
		}
	}

	return loc
}
