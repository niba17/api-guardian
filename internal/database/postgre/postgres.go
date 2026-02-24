package postgre

import (
	"io"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 🚀 BUAT INTERFACE: Agar app.go bisa memegang 'kontrak'
type DBInstance interface {
	GetDB() *gorm.DB
	io.Closer // Wajib punya fungsi Close()
}

// 🚀 STRUCT PRIVATE: Implementasi detail yang disembunyikan
type postgresDB struct {
	db *gorm.DB
}

func (p *postgresDB) GetDB() *gorm.DB {
	return p.db
}

func (p *postgresDB) Close() error {
	sqlDB, _ := p.db.DB()
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}

// InitDB sekarang mengembalikan Interface DBInstance
func InitDB(dsn string) DBInstance {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to PostgreSQL")
		return nil
	}

	log.Info().Msg("Successfully connected to PostgreSQL")
	return &postgresDB{db: db}
}
