package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

// Kita buat struct agar bisa memegang koneksi Redis
type DashboardHandler struct {
	Redis *redis.Client
}

func NewDashboardHandler(rdb *redis.Client) *DashboardHandler {
	return &DashboardHandler{Redis: rdb}
}

// --- HANDLER 1: Statistik Ringkas (Dari Redis) ---
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	ctx := context.Background()

	// Ambil data dari Redis (Otomatis 0 jika kunci tidak ada)
	total, _ := h.Redis.Get(ctx, "stats:total_requests").Int64()
	blocked, _ := h.Redis.Get(ctx, "stats:blocked_requests").Int64()
	uniqueIPs, _ := h.Redis.SCard(ctx, "stats:unique_ips").Result()

	stats := map[string]interface{}{
		"total_requests":   total,
		"blocked_requests": blocked,
		"unique_ips":       uniqueIPs,
		"avg_latency":      "15ms", // Ini bisa dihitung dinamis nanti
	}

	json.NewEncoder(w).Encode(stats)
}

// --- HANDLER 2: Log Terbaru (Dari Redis List) ---
func (h *DashboardHandler) GetRecentLogs(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	ctx := context.Background()

	// Ambil 50 log terbaru dari Redis List "stats:recent_logs"
	val, err := h.Redis.LRange(ctx, "stats:recent_logs", 0, 49).Result()

	var logs []interface{}
	for _, item := range val {
		var l interface{}
		json.Unmarshal([]byte(item), &l)
		logs = append(logs, l)
	}

	if err != nil || logs == nil {
		logs = []interface{}{}
	}

	json.NewEncoder(w).Encode(logs)
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}
