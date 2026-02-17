package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// APIKeyValidator hanya fokus pada validasi kunci akses
func APIKeyValidator(validKeys []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey := r.Header.Get("X-API-KEY")

		if clientKey == "" {
			// 👇 Tambahkan Header untuk Logger
			w.Header().Set("X-Guardian-WAF-Reason", "Missing API Key")

			log.Warn().Str("module", "auth").Str("ip", r.RemoteAddr).Msg("Missing API Key")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized", "message": "Missing API Key"}`))
			return
		}

		isValid := false
		for _, k := range validKeys {
			if k != "" && clientKey == k {
				isValid = true
				break
			}
		}

		if !isValid {
			// 👇 Tambahkan Header untuk Logger
			w.Header().Set("X-Guardian-WAF-Reason", "Invalid API Key Attempt")

			log.Warn().Str("module", "auth").Str("ip", r.RemoteAddr).Msg("Invalid API Key")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Forbidden", "message": "Invalid API Key"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
