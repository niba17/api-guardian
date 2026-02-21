package postgre

import (
	"api-guardian/internal/domain/user" // 👈 Sesuaikan dengan path model User Bos

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	var count int64
	// 1. Cek dulu apakah admin sudah ada, biar nggak duplikat
	db.Model(&user.User{}).Where("username = ?", "admin").Count(&count)

	if count > 0 {
		log.Info().Msg("✅ Admin user already exists, skipping seed.")
		return nil
	}

	// 2. Hash password "password123" sesuai permintaan Bos
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Suntikkan ke database
	admin := user.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Info().Msg("🚀 Successfully injected admin:password123 into database!")
	return nil
}
