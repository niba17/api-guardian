package dto

import "time"

// DashboardStatsResponse adalah balasan untuk widget angka di atas UI
type DashboardStatsResponse struct {
	TotalRequests int            `json:"total_requests"`
	TotalBlocked  int            `json:"total_blocked"`
	TotalSuccess  int            `json:"total_success"`
	UniqueIPs     int            `json:"unique_ips"`
	AvgLatency    int64          `json:"avg_latency_ms"` // Tambahan _ms agar UI tahu ini milidetik
	TopIPs        map[string]int `json:"top_ips"`
}

// RecentLogResponse adalah balasan untuk tabel log (versi langsing)
type RecentLogResponse struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	IP            string    `json:"ip"`
	Country       string    `json:"country"`
	City          string    `json:"city"` // 🚀 Tambahan
	Method        string    `json:"method"`
	Path          string    `json:"path"`
	Status        int       `json:"status"`
	IsBlocked     bool      `json:"is_blocked"`
	Latency       int64     `json:"latency_ms"`
	ThreatType    string    `json:"threat_type"`    // 🚀 Tambahan
	ThreatDetails string    `json:"threat_details"` // 🚀 INI DIA BINTANG UTAMANYA (Reason)!
	Browser       string    `json:"browser"`        // 🚀 Tambahan
	OS            string    `json:"os"`             // 🚀 Tambahan
	IsBot         bool      `json:"is_bot"`         // 🚀 Tambahan
}
