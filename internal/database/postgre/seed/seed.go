package seed

import (
	"api-guardian/internal/domain/user" // 👈 Sesuaikan dengan path model User Bos
	"api-guardian/pkg/hashutil"         // 🚀 IMPORT UTILITY BOS DI SINI

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Run mengeksekusi pengisian data awal ke database
func Run(db *gorm.DB) error {
	var count int64

	// 1. Cek dulu apakah admin sudah ada
	db.Model(&user.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		log.Info().Msg("✅ Admin user already exists, skipping seed.")
		return nil
	}

	// 2. Hash password "password123" menggunakan Pkg Hashutil 🚀
	hashedPassword, err := hashutil.HashPassword("password123")
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to hash seed password")
		return err
	}

	// 3. Suntikkan ke database
	admin := user.User{
		Username: "admin",
		// 🚀 Karena HashPassword sudah mengembalikan string, tidak perlu di-cast string() lagi!
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Error().Err(err).Msg("❌ Failed to insert seed admin")
		return err
	}

	log.Info().Msg("🚀 Successfully injected admin:password123 into database!")
	return nil
}
