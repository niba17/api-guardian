package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// --- KONFIGURASI PII MASKING ---

// --- STRUKTUR DATA ---
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
	// 👇 KOLOM BARU: BODY (Tipe Text biar muat banyak)
	Body string `gorm:"type:text" json:"body"`
}

// internal/middleware/logger.go

func GetIP(r *http.Request) string {
	// 1. Cek header dari Proxy Docker (X-Forwarded-For)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Ambil IP pertama sebelum koma
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// 2. Cek Real-IP
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}

	// 3. Fallback terakhir
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Jika IP adalah loopback Docker, coba tandai agar kita tahu ini akses lokal
	if ip == "172.18.0.1" {
		return "windows-host-local" // Atau IP asli Windows Bos jika tahu
	}
	return ip
}

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

// --- MIDDLEWARE UTAMA ---
func AuditLogger(geoDB *geoip2.Reader, rdb *redis.Client, db *gorm.DB, next http.Handler) http.Handler {

	// Auto Migrate akan menambahkan kolom 'body' secara otomatis
	if db != nil {
		if err := db.AutoMigrate(&SecurityLog{}); err != nil {
			log.Error().Err(err).Msg("Gagal migrasi database")
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := uuid.New().String()

		// --- 1. BACA BODY (Hati-hati di sini!) ---
		var bodyString string
		// Hanya baca body jika method POST/PUT/PATCH
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			// Baca seluruh body ke byte array
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				// Kembalikan body ke request agar bisa dibaca lagi oleh Handler berikutnya (PENTING!)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Simpan ke string dan SENSOR!
				bodyString = MaskPII(string(bodyBytes))
			}
		}

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		clientIP := GetIP(r)

		// LOGIC GEOIP: Gunakan IP dummy jika akses lokal (karena IP lokal gak ada petanya)
		targetIP := clientIP
		if clientIP == "127.0.0.1" || strings.HasPrefix(clientIP, "172.") {
			targetIP = "180.252.173.1" // Dummy IP Jakarta
		}

		countryName := "Unknown"
		cityName := "Unknown"
		if geoDB != nil {
			if ipParsed := net.ParseIP(targetIP); ipParsed != nil {
				if record, err := geoDB.City(ipParsed); err == nil {
					countryName = record.Country.Names["en"]
					cityName = record.City.Names["en"]
				}
			}
		}

		uaString := r.UserAgent()
		browser := "Other"
		os := "Other"
		isBot := false

		// Logic Deteksi Sederhana (Bisa ditingkatkan pakai library 'ua-parser')
		lowUA := strings.ToLower(uaString)
		if strings.Contains(lowUA, "chrome") {
			browser = "Chrome"
		}
		if strings.Contains(lowUA, "firefox") {
			browser = "Firefox"
		}
		if strings.Contains(lowUA, "safari") && !strings.Contains(lowUA, "chrome") {
			browser = "Safari"
		}

		if strings.Contains(lowUA, "windows") {
			os = "Windows"
		}
		if strings.Contains(lowUA, "macintosh") {
			os = "MacOS"
		}
		if strings.Contains(lowUA, "android") {
			os = "Android"
		}
		if strings.Contains(lowUA, "iphone") {
			os = "iOS"
		}

		// Deteksi Bot
		if strings.Contains(lowUA, "bot") || strings.Contains(lowUA, "spider") || strings.Contains(lowUA, "crawler") {
			isBot = true
		}

		// --- THREAT ANALYSIS (Update Jalur Header) ---
		isBlocked := wrapped.statusCode >= 400
		threatType := "None"
		threatDetails := "-"

		if isBlocked {
			// 1. Cek apakah ada pesan spesifik dari WAF di header
			wafReason := wrapped.Header().Get("X-Guardian-WAF-Reason")

			if wafReason != "" {
				threatType = "WAF"
				threatDetails = wafReason
			} else {
				// 2. Fallback logic lama Bos
				switch wrapped.statusCode {
				case 429:
					threatType = "Rate Limit"
				case 401:
					threatType = "Auth"
				case 403:
					threatType = "Access Denied"
				case 503:
					threatType = "Circuit Breaker"
				default:
					threatType = "SystemError"
				}
			}
		}

		// --- SIMPAN KE DB ---
		logData := SecurityLog{
			ID:        reqID,
			Timestamp: time.Now(),
			IP:        clientIP,
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    wrapped.statusCode,
			Latency:   duration.Milliseconds(),
			Country:   countryName,
			City:      cityName, UserAgent: uaString,
			Browser:       browser,
			OS:            os,
			IsBot:         isBot,
			IsBlocked:     isBlocked,
			ThreatType:    threatType,
			ThreatDetails: threatDetails, // 👈 Simpan detailnya di sini
			Body:          bodyString,
		}

		go func() {
			if db != nil {
				if err := db.Create(&logData).Error; err != nil {
					log.Error().Err(err).Msg("Gagal simpan log ke Database")
				}
			}
		}()

		// --- LOG KE CONSOLE (Biar Bos bisa lihat langsung) ---
		logger := log.Info()
		if isBlocked {
			logger = log.Warn()
		}

		logger.
			Interface("details", logData).
			Msg("Audit Log Entry")
	})
}
