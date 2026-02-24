package migrate

import (
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/domain/user"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Run mengeksekusi AutoMigrate GORM
func Run(db *gorm.DB) {
	log.Info().Msg("⏳ Running PostgreSQL AutoMigration...")

	err := db.AutoMigrate(
		&user.User{},
		&security_log.SecurityLog{},
	)

	if err != nil {
		log.Error().Err(err).Msg("❌ AutoMigration failed")
		return
	}

	log.Info().Msg("✅ AutoMigration completed successfully")
}
