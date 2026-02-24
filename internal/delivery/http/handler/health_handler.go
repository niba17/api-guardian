package handler

import (
	"api-guardian/internal/domain/health/interfaces"
	"api-guardian/pkg/response"
	"net/http"
)

type HealthHandler struct {
	HealthUC interfaces.HealthUsecase
}

func NewHealthHandler(uc interfaces.HealthUsecase) *HealthHandler {
	return &HealthHandler{
		HealthUC: uc,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// 1. Panggil Usecase (Biarkan Usecase yang mikir logika Redis dll)
	status := h.HealthUC.CheckHealth(r.Context())

	// 2. Balas pakai Response Wrapper standar Bos
	response.JSON(w, http.StatusOK, status)
}
