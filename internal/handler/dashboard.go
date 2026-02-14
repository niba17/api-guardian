package handler

import (
	"api-guardian/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewDashboardHandler(rdb *redis.Client, db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		Redis: rdb,
		DB:    db,
	}
}

// --- HANDLER 1: Statistik Ringkas (SEKARANG DARI POSTGRES) ---
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)

	var totalReq int64
	var blockedReq int64
	var uniqueIPs int64
	var avgLatency float64

	// 1. Hitung Total Requests (SELECT COUNT(*) FROM security_logs)
	h.DB.Model(&middleware.SecurityLog{}).Count(&totalReq)

	// 2. Hitung Blocked Requests (WHERE is_blocked = true)
	h.DB.Model(&middleware.SecurityLog{}).Where("is_blocked = ?", true).Count(&blockedReq)

	// 3. Hitung Unique IPs (SELECT COUNT(DISTINCT ip))
	h.DB.Model(&middleware.SecurityLog{}).Distinct("ip").Count(&uniqueIPs)

	// 4. Hitung Rata-rata Latency (SELECT AVG(latency))
	// Result scan ke float64, handle kalau null (data kosong)
	h.DB.Model(&middleware.SecurityLog{}).Select("COALESCE(AVG(latency), 0)").Scan(&avgLatency)

	stats := map[string]interface{}{
		"total_requests":   totalReq,
		"blocked_requests": blockedReq,
		"unique_ips":       uniqueIPs,
		"avg_latency":      int64(avgLatency), // Konversi ke int biar rapi (misal: 15ms)
	}

	json.NewEncoder(w).Encode(stats)
}

// --- HANDLER 2: Log Terbaru (Tetap Sama) ---
func (h *DashboardHandler) GetRecentLogs(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)

	var logs []middleware.SecurityLog

	// Ambil 100 log terakhir
	result := h.DB.Order("timestamp desc").Limit(100).Find(&logs)

	if result.Error != nil {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	json.NewEncoder(w).Encode(logs)
}

// --- UTILS ---
func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}
