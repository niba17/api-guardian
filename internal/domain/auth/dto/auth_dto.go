package dto

import "api-guardian/internal/domain/user"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  user.User `json:"user"`
}
