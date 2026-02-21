package handler

import (
	"api-guardian/internal/domain/dashboard/interfaces"
	"api-guardian/pkg/response"
	"net/http"
)

type DashboardHandler struct {
	DashboardUC interfaces.DashboardUsecase
}

func NewDashboardHandler(uc interfaces.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		DashboardUC: uc,
	}
}

// --- HANDLER 1: Statistik Ringkas ---
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {

	stats, err := h.DashboardUC.GetStats(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

// --- HANDLER 2: Log Terbaru ---
func (h *DashboardHandler) GetRecentLogs(w http.ResponseWriter, r *http.Request) {

	logs, err := h.DashboardUC.GetRecentLogs(r.Context(), 100)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, logs)
}
