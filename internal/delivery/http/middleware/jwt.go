package middleware

import (
	"api-guardian/pkg/jwtutil"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

// 🚀 NAMA FUNGSI DIGANTI: Menjadi JWTMiddleware agar lebih spesifik
// JWTMiddleware memeriksa apakah request memiliki token JWT yang valid
func JWTMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Unauthorized", "message": "Missing or invalid token format"}`))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtutil.ValidateToken(tokenStr, jwtKey)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Unauthorized", "message": "Invalid or expired token"}`))
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
