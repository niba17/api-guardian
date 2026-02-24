package interfaces

import "api-guardian/internal/domain/user"

// UserRepository hanya fokus pada kontrak akses data User ke database
type UserRepository interface {
	GetUserByUsername(username string) (*user.User, error)
	CreateUser(user *user.User) error
}
