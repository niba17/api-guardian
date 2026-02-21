package postgre

import (
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB membuka koneksi ke PostgreSQL
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to PostgreSQL")
		return nil
	}

	log.Info().Msg("Successfully connected to PostgreSQL")

	return db
}
