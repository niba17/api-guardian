package security_log

import "time"

type SecurityLog struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	Timestamp     time.Time `gorm:"index" json:"timestamp"`
	IP            string    `gorm:"size:50" json:"ip"`
	Method        string    `gorm:"size:10" json:"method"`
	Path          string    `json:"path"`
	Status        int       `json:"status"`
	Latency       int64     `json:"latency"`
	Country       string    `gorm:"size:100" json:"country"`
	City          string    `gorm:"size:100" json:"city"`
	UserAgent     string    `json:"user_agent"`
	Browser       string    `gorm:"size:50" json:"browser"`
	OS            string    `gorm:"size:50" json:"os"`
	IsBot         bool      `json:"is_bot"`
	IsBlocked     bool      `gorm:"index" json:"is_blocked"`
	ThreatType    string    `gorm:"size:50" json:"threat_type"`
	ThreatDetails string    `gorm:"size:255" json:"threat_details"`
	Body          string    `gorm:"type:text" json:"body"`
}

// DashboardStats dipindah ke sini (sebagai Domain Value Object)
type DashboardStats struct {
	TotalRequests int            `json:"total_requests"`
	TotalBlocked  int            `json:"total_blocked"`
	TotalSuccess  int            `json:"total_success"`
	UniqueIPs     int            `json:"unique_ips"`
	AvgLatency    int64          `json:"avg_latency"`
	TopIPs        map[string]int `json:"top_ips"`
}
