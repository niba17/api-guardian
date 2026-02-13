package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var serverHostname, _ = os.Hostname()

// --- STRUKTUR DATA (Shared dengan Handler) ---

type SecurityLog struct {
	ID         string `json:"id"`
	Timestamp  string `json:"timestamp"`
	IP         string `json:"ip"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	Latency    int64  `json:"latency"`
	Country    string `json:"country"`
	City       string `json:"city"`
	IsBlocked  bool   `json:"is_blocked"`
	ThreatType string `json:"threat_type,omitempty"`
}

type DashboardStats struct {
	TotalRequests   int64  `json:"total_requests"`
	BlockedRequests int64  `json:"blocked_requests"`
	UniqueIPs       int    `json:"unique_ips"`
	AvgLatency      string `json:"avg_latency"`
}

// --- MEMORY STORE (Internal) ---
var (
	statsMutex sync.RWMutex
	totalReq   int64
	blockedReq int64
	uniqueIPs  = make(map[string]bool)
	recentLogs []SecurityLog
)

// --- PUBLIC ACCESSORS (Untuk Handler) ---

func GetDashboardStats() DashboardStats {
	statsMutex.RLock()
	defer statsMutex.RUnlock()

	return DashboardStats{
		TotalRequests:   totalReq,
		BlockedRequests: blockedReq,
		UniqueIPs:       len(uniqueIPs),
		AvgLatency:      "15ms",
	}
}

func GetRecentLogs() []SecurityLog {
	statsMutex.RLock()
	defer statsMutex.RUnlock()

	// Copy slice agar aman
	logsCopy := make([]SecurityLog, len(recentLogs))
	copy(logsCopy, recentLogs)
	return logsCopy
}

// --- INTERNAL HELPERS ---

func saveToMemory(logEntry SecurityLog) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	// Update Stats
	totalReq++
	if logEntry.IsBlocked {
		blockedReq++
	}
	uniqueIPs[logEntry.IP] = true

	// Update Logs (Prepend)
	recentLogs = append([]SecurityLog{logEntry}, recentLogs...)
	// Keep only last 50
	if len(recentLogs) > 50 {
		recentLogs = recentLogs[:50]
	}
}

func GetIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	// Langsung sikat saja tanpa 'if', Go yang urus sisanya
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")
	host = strings.TrimPrefix(host, "::ffff:")

	return host
}

// --- MIDDLEWARE UTAMA ---

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func AuditLogger(geoDB *geoip2.Reader, rdb *redis.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := uuid.New().String()
		w.Header().Set("X-Request-ID", reqID)

		var reqBody []byte
		var reqSize int64
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			reqSize = int64(len(reqBody))
		}

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		clientIP := GetIP(r)

		// GeoIP Logic
		targetIP := clientIP
		if clientIP == "::1" || clientIP == "127.0.0.1" || strings.HasPrefix(clientIP, "172.") {
			targetIP = "180.252.173.1" // Telkomsel Dummy
		}

		countryName := "Unknown"
		cityName := "Unknown"
		if geoDB != nil {
			ipParsed := net.ParseIP(targetIP)
			if ipParsed != nil {
				if record, err := geoDB.City(ipParsed); err == nil {
					countryName = record.Country.Names["en"]
					cityName = record.City.Names["en"]
				}
			}
		}

		// Security Logic - Lebih rapi pakai switch
		isBlocked := wrapped.statusCode >= 400
		threatType := "None"

		switch wrapped.statusCode {
		case http.StatusTooManyRequests:
			threatType = "Rate Limit"
		case http.StatusForbidden:
			threatType = "Access Denied"
		case http.StatusBadRequest:
			threatType = "WAF Block"
		}

		// 1. Simpan ke Memori Dashboard
		logData := SecurityLog{
			ID:         reqID,
			Timestamp:  time.Now().Format(time.RFC3339),
			IP:         clientIP,
			Method:     r.Method,
			Path:       r.URL.Path,
			Status:     wrapped.statusCode,
			Latency:    duration.Milliseconds(),
			Country:    countryName,
			City:       cityName,
			IsBlocked:  isBlocked,
			ThreatType: threatType,
		}
		go saveToMemory(logData)

		// --- 🚀 PROSES ANTI-AMNESIA (REDIS) ---
		// --- 🚀 PROSES ANTI-AMNESIA (REDIS) ---
		go func(data SecurityLog, client *redis.Client) {
			if client == nil {
				return
			}

			// 1. Cukup buat SATU context saja untuk seluruh proses ini
			// Gunakan 10 detik agar lebih tahan banting di Windows
			bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// 2. Cek koneksi
			if err := client.Ping(bgCtx).Err(); err != nil {
				log.Error().Err(err).Msg("⚠️ Redis tidak merespon saat update stats")
				return
			}

			// 3. Jalankan Pipeline
			pipe := client.Pipeline()

			pipe.Incr(bgCtx, "stats:total_requests")
			if data.IsBlocked {
				pipe.Incr(bgCtx, "stats:blocked_requests")
			}
			pipe.SAdd(bgCtx, "stats:unique_ips", data.IP)

			logJSON, _ := json.Marshal(data)
			pipe.LPush(bgCtx, "stats:recent_logs", logJSON)
			pipe.LTrim(bgCtx, "stats:recent_logs", 0, 49)

			// 4. Eksekusi
			if _, err := pipe.Exec(bgCtx); err != nil {
				log.Error().Err(err).Msg("❌ Gagal eksekusi pipeline Redis")
			}
		}(logData, rdb)

		// 2. Log ke Console/File (Zerolog)
		// MaskPII dipanggil dari privacy.go (satu package)
		maskedReq := MaskPII(string(reqBody))
		maskedRes := MaskPII(wrapped.body.String())

		log.Info().
			Str("req_id", reqID).
			Str("ip", clientIP).
			Str("method", r.Method).
			Int("status", wrapped.statusCode).
			Int64("latency_ms", duration.Milliseconds()).
			Int64("req_size", reqSize). // Fix: reqSize sekarang dipakai
			Bool("blocked", isBlocked).
			Interface("req_body", maskedReq).
			Interface("res_body", maskedRes).
			Msg("Audit Log")
	})
}
