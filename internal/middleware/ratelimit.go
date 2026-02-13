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
	BAN_DURATION     = 1 * time.Hour
)

func RateLimiter(store storage.LimiterStore, whitelist []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip := GetIP(r) // Pastikan GetIP dari logger.go sudah bersih

		// [DEBUG] Cek IP yang masuk
		// log.Debug().Str("ip", ip).Msg("🔍 RateLimiter: Checking Incoming Request")

		// 1. CEK WHITELIST (Jalur VVIP)
		for _, allowed := range whitelist {
			// Kita buang logika "["+ip+"]" karena GetIP sudah membersihkannya
			if allowed == ip {
				log.Info().Str("ip", ip).Msg("🛡️ RateLimiter: VVIP WHITELISTED (Bypass)")
				next.ServeHTTP(w, r)
				return
			}
		}

		// 2. CEK BLACKLIST (Jalur Penjara)
		banKey := "blacklist:" + ip
		exists, err := store.Exists(ctx, banKey)
		if err != nil {
			log.Error().Err(err).Msg("⚠️ RateLimiter: Redis Error Check Blacklist")
			// Lanjut dulu kalau error cek ban, nanti dicek di limit
		}

		if exists > 0 {
			log.Warn().Str("ip", ip).Msg("⛔ RateLimiter: BLOCKED (User is Banned)")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "IP Anda di-BANNED selama 1 jam."}`))
			return
		}

		// 3. RATE LIMITING (Jalur Hitung)
		limitKey := "rate_limit:" + ip
		count, err := store.Incr(ctx, limitKey)

		// [DEBUG] Tangkap jika Redis Error/Mati
		if err != nil {
			log.Error().Err(err).Msg("💥 RateLimiter: REDIS ERROR on INCR (Fail Open)")
			// Di sini lubangnya! Kalau Redis error, dia lolos.
			// Untuk debugging, kita biarkan lolos tapi LOG-nya merah.
			next.ServeHTTP(w, r)
			return
		}

		// [DEBUG] Tampilkan counter saat ini
		// log.Debug().Str("ip", ip).Int64("hits", count).Msg("📊 RateLimiter Counter")

		// Set Expire untuk request pertama
		if count == 1 {
			store.Expire(ctx, limitKey, 1*time.Minute)
		}

		// 4. LOGIKA PELANGGARAN (Jalur Hukuman)
		if count > LIMIT_PER_MINUTE {
			violationKey := "violation:" + ip
			violationCount, _ := store.Incr(ctx, violationKey)

			if violationCount == 1 {
				store.Expire(ctx, violationKey, 10*time.Minute)
			}

			remainingLives := int64(MAX_VIOLATIONS) - violationCount

			log.Warn().Str("ip", ip).
				Int64("hits", count).
				Int64("violations", violationCount).
				Msg("🚫 RateLimiter: OVER LIMIT")

			if violationCount >= MAX_VIOLATIONS {
				log.Error().Str("ip", ip).Msg("🔥 RateLimiter: BAN HAMMER EXECUTED!")

				// Tulis Blacklist
				err := store.Set(ctx, banKey, "banned", BAN_DURATION)
				if err != nil {
					log.Error().Err(err).Msg("❌ RateLimiter: Gagal Tulis Blacklist ke Redis")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "Selamat! Anda resmi di-BANNED."}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf(`{"error": "Terlalu Cepat! Sisa nyawa: %d"}`, remainingLives)))
			return
		}

		// Kalau aman, lanjut
		next.ServeHTTP(w, r)
	})
}
