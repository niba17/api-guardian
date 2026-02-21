package middleware

import (
	"api-guardian/internal/domain/security_log"
	"api-guardian/internal/usecase"
	"api-guardian/pkg/uaparser"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

func Logger(geoDB *geoip2.Reader, logUC *usecase.LogUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Capture & Mask Request Body
		var bodyStr string
		if r.Method != http.MethodGet && r.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 2048))
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// Memanggil MaskPII yang ada di file mask.go (satu package)
			bodyStr = MaskPII(string(bodyBytes))
		}

		// 🚀 Bungkus writer asli dengan responseWriter kita
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 0}

		// 2. Jalankan Handler berikutnya
		next.ServeHTTP(wrapped, r)

		// Fallback jika status code tidak terisi
		if wrapped.statusCode == 0 {
			wrapped.statusCode = http.StatusOK
		}

		// 🚀 DEKLARASI VARIABEL (Pastikan semua nama ini dipanggil di bawah)
		latency := time.Since(start)
		wafReason := wrapped.Header().Get("X-Guardian-WAF-Reason")
		reqIP := GetIP(r)
		reqUA := r.UserAgent()
		reqPath := r.URL.Path
		reqMethod := r.Method
		logID := uuid.New().String()
		browser, os, isBot := uaparser.Parse(reqUA)

		// 📺 A. CETAK KE TERMINAL (Memakai reqMethod, reqPath, dll)
		statusColor := "\033[32m"
		if wrapped.statusCode >= 400 {
			statusColor = "\033[33m"
		}
		if wrapped.statusCode >= 500 {
			statusColor = "\033[31m"
		}
		fmt.Printf("📝 %s %s%d%s | %10s | %15s | %s %s\n",
			time.Now().Format("15:04:05"), statusColor, wrapped.statusCode, "\033[0m",
			latency, reqIP, reqMethod, reqPath)

		// 📁 B. TULIS KE app.log (Memakai semua variabel nganggur tadi)
		logEvent := log.Info()
		if wrapped.statusCode >= 400 {
			logEvent = log.Warn()
		}

		logEvent.
			Str("log_id", logID).
			Str("ip", reqIP).
			Str("method", reqMethod).
			Str("path", reqPath).
			Int("status", wrapped.statusCode).
			Int64("latency_ms", latency.Milliseconds()).
			Str("browser", browser).
			Str("os", os).
			Bool("is_bot", isBot).
			Str("body", bodyStr).
			Msg("AccessLog")

		// 3. Offload ke LogUsecase (Asynchronous)
		go func(status int, lat time.Duration, path, method, ua, ip, body, reason string) {
			browser, os, isBot := uaparser.Parse(ua)

			logData := &security_log.SecurityLog{
				ID:        uuid.New().String(),
				Timestamp: time.Now().UTC(),
				IP:        ip,
				Method:    method,
				Path:      path,
				Status:    status,
				Latency:   lat.Milliseconds(),
				UserAgent: ua,
				Browser:   browser,
				OS:        os,
				IsBot:     isBot,
				Body:      body,
			}

			enrichLogData(logData, reason, geoDB)
			logUC.LogAndEvaluate(context.Background(), logData, reason)
		}(wrapped.statusCode, latency, reqPath, reqMethod, reqUA, reqIP, bodyStr, wafReason)
	})
}

// --- HELPER STRUCTS & FUNCTIONS ---

// responseWriter adalah interceptor untuk menangkap status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.statusCode == 0 {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

func enrichLogData(l *security_log.SecurityLog, reason string, geo *geoip2.Reader) {
	l.IsBlocked = l.Status >= 400
	l.ThreatType = "None"
	l.ThreatDetails = "-"

	if geo != nil {
		ip := net.ParseIP(l.IP)
		if record, err := geo.City(ip); err == nil {
			l.Country = record.Country.Names["en"]
			l.City = record.City.Names["en"]
		}
	}

	if l.IsBlocked {
		if reason != "" {
			l.ThreatType = "WAF"
			l.ThreatDetails = reason
		} else if l.Status == 429 {
			l.ThreatType = "RateLimit"
			l.ThreatDetails = "Too Many Requests"
		} else {
			l.ThreatType = "ResponseError"
		}
	}
}
