package handler

import (
	"api-guardian/internal/delivery/http/dto"
	"api-guardian/internal/usecase"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	AuthUC usecase.AuthUsecase // 👈 Dependensinya ke Usecase, BUKAN Repo
}

// Constructor sekarang menerima Usecase
func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUC: uc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Setup CORS sederhana (Bisa dihapus kalau sudah ada middleware CORS global)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode Request ke DTO
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. Panggil Usecase (Otak Bisnis)
	resp, err := h.AuthUC.Login(req)
	if err != nil {
		// Kita samarkan errornya jadi "Unauthorized" biar aman
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 3. Kirim Response Sukses
	json.NewEncoder(w).Encode(resp)
}
