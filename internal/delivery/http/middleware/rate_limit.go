package middleware

import (
	rlInterfaces "api-guardian/internal/domain/rate_limit/interfaces"
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func RateLimit(rdb rlInterfaces.RateLimitRepository, banUC rlInterfaces.BanUsecase, whitelist []string, refillRate float64, burstCapacity int, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip := GetIP(r)

		// 1. CEK BLACKLIST
		banKey := "blacklist:" + ip
		if isBanned, _ := rdb.Exists(ctx, banKey); isBanned > 0 {
			// 🛑 log.Warn() DIHAPUS: Biarkan logger.go yang mencatat log 403-nya biar rapi!
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
		allowed, remaining, err := rdb.TakeToken(ctx, limitKey, 1, burstCapacity, refillRate)

		if err != nil {
			// 🚀 UBAH KE DEBUG: Tidak akan tampil di terminal secara default,
			// tapi WAF Logger (logger.go) tetap akan mencatat HTTP 503-nya!
			log.Debug().Err(err).Msg("REDIS ERROR: Fail-Closed triggered")

			// Tambahkan WAF Reason agar UI Dashboard tahu ini error karena sistem keamanan offline
			w.Header().Set("X-Guardian-WAF-Reason", "Fail-Closed: WAF Offline")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		// 4. JIKA LIMIT HABIS
		if !allowed {
			w.Header().Set("X-Guardian-WAF-Reason", "Rate Limit Exceeded")
			go banUC.ExecuteAutoBan(context.Background(), ip)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too Many Requests. Slow down!"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
