package model

// LogEntry merepresentasikan satu baris di app.log
type LogEntry struct {
	Level     string `json:"level"`
	Time      string `json:"time"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Duration  string `json:"duration"`
	UserAgent string `json:"user_agent"`
	Error     string `json:"error,omitempty"`
	Reason    string `json:"reason,omitempty"` // Kalau diblokir WAF
}

// DashboardStats adalah rangkuman siap saji untuk Frontend
type DashboardStats struct {
	TotalRequests int            `json:"total_requests"`
	TotalBlocked  int            `json:"total_blocked"` // Status 4xx/5xx
	TotalSuccess  int            `json:"total_success"` // Status 2xx
	TopIPs        map[string]int `json:"top_ips"`       // Siapa pengunjung paling rajin
	Attacks       int            `json:"total_attacks"` // Terdeteksi WAF
}
