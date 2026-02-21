package handler

import (
	"api-guardian/internal/domain/auth/dto"
	"api-guardian/internal/domain/auth/interfaces"
	"api-guardian/pkg/response"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	AuthUC interfaces.AuthUsecase
}

func NewAuthHandler(uc interfaces.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUC: uc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// 1. Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, interfaces.ErrInvalidJSONFormat)
		return
	}

	// 2. Panggil Usecase
	resp, err := h.AuthUC.Login(req)
	if err != nil {
		response.Error(w, err)
		return
	}

	// 3. Sukses
	response.JSON(w, http.StatusOK, resp)
}
