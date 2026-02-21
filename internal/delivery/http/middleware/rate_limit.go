package middleware

import (
	"api-guardian/internal/domain/rate_limit/interfaces"
	"api-guardian/internal/usecase" // 👈 Import Usecase
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// RateLimit sekarang menerima BanUsecase!
func RateLimit(store interfaces.LimitRepository, banUC *usecase.BanUsecase, whitelist []string, refillRate float64, burstCapacity int, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip := GetIP(r) // 👈 Ingat, fungsi ini masih misteri! 😁

		// 1. CEK BLACKLIST
		banKey := "blacklist:" + ip
		if isBanned, _ := store.Exists(ctx, banKey); isBanned > 0 {
			log.Warn().Str("ip", ip).Msg("Refusing request from Blacklisted IP")
			w.Header().Set("X-Guardian-WAF-Reason", "IP Blacklisted")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Access Denied", "message": "Your IP is permanently banned"}`))
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
		// 🚀 GUNAKAN PARAMETER DINAMIS DI SINI
		allowed, remaining, err := store.TakeToken(ctx, limitKey, 1, burstCapacity, refillRate)

		if err != nil {
			log.Error().Err(err).Msg("REDIS ERROR: Fail-Closed")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		// 4. JIKA LIMIT HABIS
		if !allowed {
			w.Header().Set("X-Guardian-WAF-Reason", "Rate RateLimit Exceeded")

			// 🚀 PANGGIL USECASE
			go banUC.ExecuteAutoBan(context.Background(), ip)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too Many Requests. Slow down!"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
