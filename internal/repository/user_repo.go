package repository

import (
	"api-guardian/internal/domain/user"
	"api-guardian/internal/domain/user/interfaces"
	"errors"

	"gorm.io/gorm"
)

// 🚀 STRUCT KECIL: Privat untuk package repository
type gormUserRepo struct {
	db *gorm.DB
}

// 🚀 RETURN INTERFACE: Konsisten dengan kasta domain
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &gormUserRepo{db: db}
}

func (r *gormUserRepo) GetUserByUsername(username string) (*user.User, error) {
	var u user.User
	result := r.db.Where("username = ?", username).First(&u)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *gormUserRepo) CreateUser(u *user.User) error {
	return r.db.Create(u).Error
}
