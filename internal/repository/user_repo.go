package repository

import (
	"api-guardian/internal/domain/user"
	"api-guardian/internal/domain/user/interfaces"
	"errors"

	"gorm.io/gorm"
)

// gormUserRepo dibuat private (huruf kecil) agar tidak bisa diinstansiasi langsung dari luar
type gormUserRepo struct {
	db *gorm.DB
}

// NewUserRepository sekarang mengembalikan interfaces.UserRepository
// Ini mengikuti pola SecurityLogRepository yang mengembalikan interface
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &gormUserRepo{db: db}
}

// GetUserByUsername mencari user berdasarkan username
func (r *gormUserRepo) GetUserByUsername(username string) (*user.User, error) {
	var user user.User

	result := r.db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Menggunakan error dari domain/interfaces agar sinkron dengan Usecase
			return nil, interfaces.ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// CreateUser menyimpan user baru ke database
func (r *gormUserRepo) CreateUser(user *user.User) error {
	return r.db.Create(user).Error
}
