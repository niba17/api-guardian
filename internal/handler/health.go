package handler

import (
	"api-guardian/internal/storage" // Import storage
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// HealthCheck sekarang terima interface
func HealthCheck(store storage.LimiterStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redisStatus := "Connected"

		// Cek via interface Ping
		if err := store.Ping(context.Background()); err != nil {
			redisStatus = "Disconnected: " + err.Error()
		}

		status := map[string]string{
			"system":           "Healthy",
			"redis_connection": redisStatus,
			"circuit_breaker":  "Closed (Normal)", // Hardcode dulu atau ambil dari state CB
			"time":             time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
