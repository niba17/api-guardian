package interfaces

import (
	"api-guardian/internal/domain/auth/dto"
)

// --- Usecase Interfaces ---

// AuthUsecase adalah "Otak" yang mengatur alur login dan registrasi.
type AuthUsecase interface {
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	Register(req dto.LoginRequest) error
}
