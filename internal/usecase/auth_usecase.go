package usecase

import (
	"api-guardian/internal/delivery/http/dto"
	"api-guardian/internal/repository"
	"api-guardian/pkg/hashutil"
	"api-guardian/pkg/jwtutil"
	"errors"
	"fmt"
	"strings"
)

type AuthUsecase interface {
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	Register(req dto.LoginRequest) error
}

type authUsecase struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
}

func NewAuthUsecase(repo *repository.UserRepository, key string) AuthUsecase {
	return &authUsecase{
		userRepo: repo,
		jwtKey:   []byte(key),
	}
}

// Implementasi Method Login (Pindahan logic dari Handler)
func (u *authUsecase) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	// 1. Cari User
	user, err := u.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	// Gunakan strings.TrimSpace untuk jaga-jaga
	if !hashutil.CheckPasswordHash(strings.TrimSpace(req.Password), strings.TrimSpace(user.PasswordHash)) {
		fmt.Println("DEBUG: Password mismatch!")
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	// 2. Cek Password dengan Sterilisasi
	// Kita trim input dan hash dari DB untuk membuang karakter siluman
	inputPassword := strings.TrimSpace(req.Password)
	storedHash := strings.TrimSpace(user.PasswordHash)

	if !hashutil.CheckPasswordHash(inputPassword, storedHash) {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	// 3. Generate Token
	token, err := jwtutil.GenerateToken(uint(user.ID), user.Username, user.Role, u.jwtKey)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *authUsecase) Register(req dto.LoginRequest) error {
	hashed, _ := hashutil.HashPassword(req.Password)
	// Buat objek user baru (sesuaikan dengan struct domain Bos)
	// user := &domain.User{Username: req.Username, PasswordHash: hashed, Role: "admin"}
	// return u.userRepo.CreateUser(user)

	// Untuk tes cepat, cetak saja hash yang dihasilkan aplikasi:
	fmt.Printf("HASH BARU UNTUK [%s]: %s\n", req.Password, hashed)
	return nil
}
