package auth

import (
	"api-guardian/internal/domain/auth/interfaces"
	"api-guardian/pkg/hashutil"
	"api-guardian/pkg/jwtutil"
)

// BcryptHasher adalah implementasi dari interfaces.PasswordHasher
type BcryptHasher struct{}

func NewBcryptHasher() interfaces.PasswordHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(p string) (string, error) {
	return hashutil.HashPassword(p)
}

func (h *BcryptHasher) Compare(p, hash string) bool {
	return hashutil.CheckPasswordHash(p, hash)
}

// JWTProvider adalah implementasi dari interfaces.TokenProvider
type JWTProvider struct {
	secret []byte
}

func NewJWTProvider(secret string) interfaces.TokenProvider {
	return &JWTProvider{secret: []byte(secret)}
}

func (j *JWTProvider) Generate(id uint, u, r string) (string, error) {
	return jwtutil.GenerateToken(id, u, r, j.secret)
}
