package handler

import (
	// 👈 Import Domain
	"api-guardian/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type DashboardHandler struct {
	Redis *redis.Client
	Repo  repository.SecurityLogRepository // 👈 Ganti GORM dengan Repo
}

func NewDashboardHandler(rdb *redis.Client, repo repository.SecurityLogRepository) *DashboardHandler {
	return &DashboardHandler{
		Redis: rdb,
		Repo:  repo, // 👈 Inject Repo
	}
}

// --- HANDLER 1: Statistik Ringkas ---
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	// Panggil logika dari Repo, bukan direct DB
	stats, err := h.Repo.GetStats()
	if err != nil {
		http.Error(w, "Fail to get stats", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// --- HANDLER 2: Log Terbaru ---
func (h *DashboardHandler) GetRecentLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.Repo.GetRecent(100)
	if err != nil {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	json.NewEncoder(w).Encode(logs)
}
