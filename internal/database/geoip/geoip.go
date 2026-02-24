package geoip

import (
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

// InitGeoIP membuka koneksi ke database lokal MaxMind (.mmdb)
// Jika file tidak ditemukan, tidak akan crash (mengembalikan nil)
func InitGeoIP(dbPath string) *geoip2.Reader {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Warn().
			Str("path", dbPath).
			Err(err).
			Msg("🌍 GeoIP Database NOT FOUND. Location will be 'Unknown'")
		return nil
	}

	log.Info().
		Str("path", dbPath).
		Msg("🌍 GeoIP Database Loaded Successfully")

	return db
}
