package usecase

import (
	"api-guardian/internal/domain/auth/dto"
	authInterface "api-guardian/internal/domain/auth/interfaces"
	userInterface "api-guardian/internal/domain/user/interfaces"
	"api-guardian/pkg/hashutil"
	"fmt"
	"strings"
)

type authUsecase struct {
	userRepo userInterface.UserRepository
	hasher   authInterface.PasswordHasher
	tokenSvc authInterface.TokenProvider
}

// Constructor yang "Luwes" (Menerima apa saja yang penting sesuai interface)
func NewAuthUsecase(
	repo userInterface.UserRepository,
	h authInterface.PasswordHasher,
	ts authInterface.TokenProvider,
) authInterface.AuthUsecase {
	return &authUsecase{
		userRepo: repo,
		hasher:   h,
		tokenSvc: ts,
	}
}

func (u *authUsecase) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	// 1. Tanya Repo
	user, err := u.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return dto.LoginResponse{}, authInterface.ErrInvalidCredentials
	}

	// 2. Tanya Hasher
	if !u.hasher.Compare(strings.TrimSpace(req.Password), user.PasswordHash) {
		return dto.LoginResponse{}, authInterface.ErrInvalidCredentials
	}

	// 3. Tanya Token Service
	token, err := u.tokenSvc.Generate(uint(user.ID), user.Username, user.Role)
	if err != nil {
		return dto.LoginResponse{}, authInterface.ErrInternalServer
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
