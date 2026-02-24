package dashboard

// DashboardStats adalah Entity/Value Object khusus milik domain Dashboard
type Dashboard struct {
	TotalRequests int            `json:"total_requests"`
	TotalBlocked  int            `json:"total_blocked"`
	TotalSuccess  int            `json:"total_success"`
	UniqueIPs     int            `json:"unique_ips"`
	AvgLatency    int64          `json:"avg_latency"`
	TopIPs        map[string]int `json:"top_ips"`
}
