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
			bodyStr = MaskPII(string(bodyBytes))
		}

		wrapped := &responseWriter{ResponseWriter: w, statusCode: 0}
		next.ServeHTTP(wrapped, r)

		if wrapped.statusCode == 0 {
			wrapped.statusCode = http.StatusOK
		}

		// 🚀 DATA GATHERING (Lakukan di awal agar Log & DB sinkron)
		latency := time.Since(start)
		wafReason := wrapped.Header().Get("X-Guardian-WAF-Reason")
		reqIP := GetIP(r)

		// 🛠️ MOCK IP UNTUK TESTING (Opsional: Hapus jika sudah deploy ke VPS beneran)
		if reqIP == "localhost" || reqIP == "::1" || reqIP == "127.0.0.1" {
			reqIP = "180.252.170.158"
		}

		reqUA := r.UserAgent()
		browser, os, isBot := uaparser.Parse(reqUA)
		logID := uuid.New().String()

		// 🌍 GEOIP LOOKUP (Pindahkan ke sini agar Zerolog bisa pakai)
		country, city := "Unknown", "Unknown"
		if geoDB != nil {
			ip := net.ParseIP(reqIP)
			if record, err := geoDB.City(ip); err == nil {
				if c, ok := record.Country.Names["en"]; ok {
					country = c
				}
				if c, ok := record.City.Names["en"]; ok {
					city = c
				}
			}
		}

		// 📺 A. CETAK KE TERMINAL (Fast View)
		statusColor := "\033[32m"
		if wrapped.statusCode >= 400 {
			statusColor = "\033[33m"
		}
		if wrapped.statusCode >= 500 {
			statusColor = "\033[31m"
		}

		fmt.Printf("📝 %s %s%d%s | %10s | %15s (%s) | %s %s\n",
			time.Now().Format("15:04:05"), statusColor, wrapped.statusCode, "\033[0m",
			latency, reqIP, country, r.Method, r.URL.Path)

		// 📁 B. TULIS KE app.log (Sekarang lengkap dengan Country/City!)
		logEvent := log.Info()
		if wrapped.statusCode >= 400 {
			logEvent = log.Warn()
		}

		logEvent.
			Str("log_id", logID).
			Str("ip", reqIP).
			Str("country", country). // 🌍 Muncul di log!
			Str("city", city).       // 🌍 Muncul di log!
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.statusCode).
			Int64("latency_ms", latency.Milliseconds()).
			Str("browser", browser).
			Str("os", os).
			Bool("is_bot", isBot).
			Msg("AccessLog")

		// 3. Offload ke LogUsecase (Database)
		go func(status int, lat time.Duration, path, method, ua, ip, body, reason, cty, ctr string) {
			logData := &security_log.SecurityLog{
				ID:        uuid.New().String(),
				Timestamp: time.Now().UTC(),
				IP:        ip,
				Country:   ctr,
				City:      cty,
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

			enrichStatusData(logData, reason)
			logUC.LogAndEvaluate(context.Background(), logData, reason)
		}(wrapped.statusCode, latency, r.URL.Path, r.Method, reqUA, reqIP, bodyStr, wafReason, city, country)
	})
}

// --- HELPERS ---

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

// Pisahkan fungsi enrich hanya untuk status, Geo sudah diproses di atas
func enrichStatusData(l *security_log.SecurityLog, reason string) {
	l.IsBlocked = l.Status >= 400
	l.ThreatType = "None"
	l.ThreatDetails = "-"

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
