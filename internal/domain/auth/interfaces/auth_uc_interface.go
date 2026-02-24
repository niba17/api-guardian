package interfaces

import (
	"api-guardian/internal/domain/auth/dto"
)

type AuthUsecase interface {
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	Register(req dto.LoginRequest) error
}
