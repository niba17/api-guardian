package repository

import (
	"api-guardian/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository menerima *gorm.DB
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByUsername mencari user menggunakan GORM
func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User

	// Karena kita sudah update struct User di domain dengan tag GORM,
	// Perintah ini akan otomatis mapping kolom database ke struct.
	result := r.db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Kita return error text biasa, atau custom error domain jika ada
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
