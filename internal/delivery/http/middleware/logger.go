package middleware

import (
	"api-guardian/internal/domain"
	"api-guardian/internal/repository"
	"api-guardian/pkg/maskutil"
	"api-guardian/pkg/uaparser"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

// AuditLogger sekarang menerima DUA Repository: LogRepo (GORM) dan BlRepo (Redis)
func AuditLogger(geoDB *geoip2.Reader, logRepo repository.SecurityLogRepository, blRepo repository.BlacklistRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// 1. Capture Request Body
		var bodyStr string
		if r.Method != http.MethodGet && r.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 2048))
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			bodyStr = maskutil.MaskPII(string(bodyBytes))
		}

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 2. Jalankan Handler (HANYA SEKALI, BOS! Saya hapus yang duplikat)
		next.ServeHTTP(wrapped, r)

		// 3. Offload Logging & Blacklisting
		go func(c context.Context, status int, latency time.Duration, path string, method string, ua string, ip string) {
			// Ambil reason dari header yang diset oleh WAF atau RateLimiter
			wafReason := wrapped.Header().Get("X-Guardian-Waf-Reason")

			browser, os, isBot := uaparser.Parse(ua)
			logData := &domain.SecurityLog{
				ID:        uuid.New().String(),
				Timestamp: time.Now().UTC(),
				IP:        ip,
				Method:    method,
				Path:      path,
				Status:    status,
				Latency:   latency.Milliseconds(),
				UserAgent: ua,
				Browser:   browser,
				OS:        os,
				IsBot:     isBot,
				Body:      bodyStr,
			}

			enrichLogData(logData, wafReason, geoDB)

			// 🔥 LOGIKA AUTO-BAN
			// Kita hanya hukum jika reason datang dari WAF (serangan aktif)
			// Dan jangan hukum lagi kalau statusnya memang sudah Banned
			if wafReason != "" && wafReason != "IP Permanently Banned" {
				vCount, _ := blRepo.IncrViolation(c, ip)

				if vCount == 1 {
					_ = blRepo.SetExpire(c, "violation:"+ip, 1*time.Hour)
				}

				if vCount >= 5 {
					log.Error().Str("ip", ip).Msg("🚨 AUTO-BAN EXECUTED")
					_ = blRepo.SetBan(c, ip, 24*time.Hour)
				}
			}

			_ = logRepo.Save(logData)
		}(ctx, wrapped.statusCode, time.Since(start), r.URL.Path, r.Method, r.UserAgent(), GetIP(r))
	})
}

// Helper untuk deteksi status & geo
func enrichLogData(l *domain.SecurityLog, reason string, geo *geoip2.Reader) {
	l.IsBlocked = l.Status >= 400
	l.ThreatType = "None"
	l.ThreatDetails = "-"

	if geo != nil {
		ip := net.ParseIP(l.IP)
		if record, err := geo.City(ip); err == nil {
			l.Country = record.Country.Names["en"] // Ambil nama negara
			l.City = record.City.Names["en"]       // Ambil nama kota
		}
	}

	if l.IsBlocked {
		if reason != "" {
			l.ThreatType = "WAF"
			l.ThreatDetails = reason
		} else {
			l.ThreatType = "ResponseError"
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
