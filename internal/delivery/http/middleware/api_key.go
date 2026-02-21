package middleware

import (
	"crypto/subtle" // 👈 Import ini untuk komparasi rahasia yang aman
	"net/http"

	"github.com/rs/zerolog/log"
)

// APIKeyValidator menggunakan pattern "Middleware Factory"
func APIKeyValidator(validKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientKey := r.Header.Get("X-API-KEY")
			ip := GetIP(r) // 🚀 Gunakan GetIP agar tembus Load Balancer!

			if clientKey == "" {
				w.Header().Set("X-Guardian-WAF-Reason", "Missing API Key")
				log.Warn().Str("module", "auth").Str("ip", ip).Msg("Missing API Key")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Unauthorized", "message": "Missing API Key"}`))
				return
			}

			isValid := false
			clientKeyBytes := []byte(clientKey) // Ubah ke byte untuk subtle.ConstantTimeCompare

			for _, k := range validKeys {
				if k != "" {
					// 🛡️ Mencegah hacker menebak API Key lewat selisih waktu eksekusi
					if subtle.ConstantTimeCompare(clientKeyBytes, []byte(k)) == 1 {
						isValid = true
						break
					}
				}
			}

			if !isValid {
				w.Header().Set("X-Guardian-WAF-Reason", "Invalid API Key Attempt")
				log.Warn().Str("module", "auth").Str("ip", ip).Msg("Invalid API Key")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "Forbidden", "message": "Invalid API Key"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
