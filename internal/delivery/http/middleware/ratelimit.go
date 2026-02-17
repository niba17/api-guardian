package middleware

import (
	"api-guardian/internal/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	LIMIT_PER_MINUTE = 20
	MAX_VIOLATIONS   = 5
	BAN_DURATION     = 24 * time.Hour
)

func RateLimiter(store storage.LimiterStore, whitelist []string, next http.Handler) http.Handler {
	// Ambil config dari env atau gunakan default
	refillRate := 0.33 // Default: 1 token tiap 3 detik
	burstCapacity := 10

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip := GetIP(r)

		// 1. CEK BLACKLIST (Prioritas Utama)
		banKey := "blacklist:" + ip
		if isBanned, _ := store.Exists(ctx, banKey); isBanned > 0 {
			log.Warn().Str("ip", ip).Msg("Refusing request from Blacklisted IP")
			w.Header().Set("X-Guardian-WAF-Reason", "IP Blacklisted") // Konsisten pake header ini
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Access Denied", "message": "Your IP is permanently banned due to security violations"}`))
			return
		}

		// 2. CEK WHITELIST
		for _, allowed := range whitelist {
			if allowed != "" && allowed == ip {
				next.ServeHTTP(w, r)
				return
			}
		}

		// 3. TOKEN BUCKET LOGIC 🛡️
		limitKey := "bucket:" + ip
		allowed, remaining, err := store.TakeToken(ctx, limitKey, 0, burstCapacity, refillRate)

		if err != nil {
			log.Error().Err(err).Msg("REDIS ERROR: Fail-Closed")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Header standar industri
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			violationKey := "violation:" + ip
			vCount, _ := store.Incr(ctx, violationKey)

			// 👇 TAMBAHKAN INI LAGI, BOS! Biar violation gak abadi
			if vCount == 1 {
				store.Expire(ctx, violationKey, 10*time.Minute)
			}

			if vCount >= MAX_VIOLATIONS {
				w.Header().Set("X-Guardian-WAF-Reason", "Ban Hammer Executed")
				log.Error().Str("ip", ip).Msg("BAN HAMMER EXECUTED!")

				// Tulis secara sinkron, jangan pakai goroutine biar pasti tersimpan
				store.Set(ctx, banKey, "banned", BAN_DURATION)

				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "You are permanently banned"}`))
				return
			}

			w.Header().Set("X-Guardian-Waf-Reason", "Rate Limit Exceeded (Token Bucket)")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
