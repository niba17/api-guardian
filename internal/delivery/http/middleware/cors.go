package middleware

import (
	"net/http"
)

// CORSMiddleware sekarang menerima daftar origin dari konfigurasi
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// 1. Cek apakah origin dari request ada di whitelist kita
			allow := false
			if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
				allow = true // Buka untuk publik (Public API)
				origin = "*" // Set origin menjadi *
			} else {
				for _, o := range allowedOrigins {
					if o == origin {
						allow = true
						break
					}
				}
			}

			// 2. Set Header jika diizinkan
			if allow && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				// Wajib jika frontend mengirimkan kredensial (seperti Cookie/Token) lintas domain
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// 3. Header Standar
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization, X-Forwarded-For")

			// 🚀 4. PENTING: Izinkan Frontend membaca header WAF dan Rate Limit kita!
			w.Header().Set("Access-Control-Expose-Headers", "X-RateLimit-Remaining, X-Guardian-WAF-Reason")

			// 5. Tangani Preflight OPTIONS
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
