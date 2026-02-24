package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	locInterfaces "api-guardian/internal/domain/location/interfaces"
	"api-guardian/internal/domain/security_log"
	logInterfaces "api-guardian/internal/domain/security_log/interfaces"
	"api-guardian/pkg/uaparser"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	recentSpamLogs sync.Map
)

func Logger(locUC locInterfaces.LocationUsecase, logUC logInterfaces.SecurityLogUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Capture & Mask Request Body
		var bodyStr string
		if r.Method != http.MethodGet && r.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 2048))
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// 🚀 Panggil langsung, karena masih 1 rumah di package middleware
			bodyStr = MaskPII(string(bodyBytes))
		}

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
			body:           bytes.NewBuffer(nil),
		}
		next.ServeHTTP(wrapped, r)

		if wrapped.statusCode == 0 {
			wrapped.statusCode = http.StatusOK
		}

		resBodyStr := wrapped.body.String()

		// 🚀 Panggil langsung!
		reqIP := GetIP(r)

		if reqIP == "localhost" || reqIP == "::1" || reqIP == "127.0.0.1" {
			reqIP = "180.252.170.158"
		}

		// 🚀 3. SMART FILTERING & ANTI-SPAM
		path := r.URL.Path

		isDashboard := strings.HasPrefix(path, "/api/dashboard")
		isHealthCheck := path == "/health" || path == "/status" || path == "/metrics"

		// Cek apakah ini IP kita sendiri (Orang Dalam)
		isInternalIP := reqIP == "180.252.170.158" || reqIP == "127.0.0.1" || reqIP == "localhost" || reqIP == "::1"

		// 🛑 A. ATURAN BISU (MUTE) UNTUK REQUEST SUKSES (< 400)
		if wrapped.statusCode < 400 {
			if isDashboard || (isHealthCheck && isInternalIP) {
				return // DIBUANG! Jangan dicatat.
			}
		}

		// 🛡️ B. Filter Anti-Spam (Debounce 2 Detik) HANYA UNTUK 403 (BANNED)
		// Biarkan 429 lolos semua agar kita bisa melihat proses hacker menggedor pintu!
		if wrapped.statusCode == 403 {
			spamKey := fmt.Sprintf("%s_%s_%d", reqIP, path, wrapped.statusCode)

			if lastSeen, exists := recentSpamLogs.Load(spamKey); exists {
				if time.Since(lastSeen.(time.Time)) < 2*time.Second {
					return // 🛑 BUANG DUPLIKAT SPAM 403!
				}
			}
			recentSpamLogs.Store(spamKey, time.Now())
		}

		// --- LANJUT CATAT KE DB ---

		latency := time.Since(start)
		wafReason := wrapped.Header().Get("X-Guardian-WAF-Reason")

		reqUA := r.UserAgent()
		browser, os, isBot := uaparser.Parse(reqUA)
		logID := uuid.New().String()

		country, city := "Unknown", "Unknown"
		if locUC != nil {
			loc := locUC.GetLocationByIP(r.Context(), reqIP)
			country = loc.Country
			city = loc.City
		}

		statusColor := "\033[32m"
		if wrapped.statusCode >= 400 {
			statusColor = "\033[33m"
		}
		if wrapped.statusCode >= 500 {
			statusColor = "\033[31m"
		}

		// Log ke Terminal
		fmt.Printf("📝 %s %s%d%s | %10s | %15s (%s) | %s %s\n",
			time.Now().Format("15:04:05"), statusColor, wrapped.statusCode, "\033[0m",
			latency, reqIP, country, r.Method, r.URL.Path)

		// Log ke File/Zerolog
		logEvent := log.Info()
		if wrapped.statusCode >= 400 {
			logEvent = log.Warn()
		}

		logEvent.
			Str("log_id", logID).
			Str("ip", reqIP).
			Str("country", country).
			Str("city", city).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.statusCode).
			Int64("latency_ms", latency.Milliseconds()).
			Str("browser", browser).
			Str("os", os).
			Bool("is_bot", isBot).
			Str("waf_reason", wafReason).
			Msg("AccessLog")

		// Log ke Database
		go func(status int, lat time.Duration, path, method, ua, ip, reqBody, reason, cty, ctr, resBody string) {
			logData := &security_log.SecurityLog{
				ID:        uuid.New().String(), // Atau logID kalau mau sinkron
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
				Body:      reqBody,
			}

			enrichStatusData(logData, reason, resBody)
			logUC.LogAndEvaluate(context.Background(), logData, reason)
		}(wrapped.statusCode, latency, r.URL.Path, r.Method, reqUA, reqIP, bodyStr, wafReason, city, country, resBodyStr)
	})
}

// --- HELPERS ---

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
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
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func enrichStatusData(l *security_log.SecurityLog, reason string, resBody string) {
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
			var errResp map[string]interface{}
			if err := json.Unmarshal([]byte(resBody), &errResp); err == nil {
				if errMsg, ok := errResp["error"].(string); ok {
					l.ThreatDetails = errMsg
					return
				}
				if msg, ok := errResp["message"].(string); ok {
					l.ThreatDetails = msg
					return
				}
			}
			if resBody != "" && len(resBody) < 200 {
				l.ThreatDetails = resBody
			} else {
				l.ThreatDetails = fmt.Sprintf("Error (Status %d)", l.Status)
			}
		}
	}
}
