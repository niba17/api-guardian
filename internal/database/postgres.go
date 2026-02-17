package database

import (
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB membuka koneksi ke PostgreSQL
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger Silent biar terminal gak penuh query SQL,
		// ganti ke logger.Info kalau mau debug query
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		// ❌ SALAH: log.Warn().Err(err).Msg("Failed... %v", err)
		// ✅ BENAR: .Err(err) sudah menangani errornya, .Msg() cukup pesannya saja.
		// Gunakan log.Error() karena koneksi DB itu krusial.
		log.Error().Err(err).Msg("Failed to connect to PostgreSQL")
		return nil // Return nil agar aplikasi tau koneksi gagal
	}

	// ✅ BENAR: Log Info untuk sukses (dan perbaiki teksnya bukan GeoIP)
	log.Info().Msg("Successfully connected to PostgreSQL")

	return db
}
