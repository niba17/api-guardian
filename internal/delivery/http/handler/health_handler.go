package handler

import (
	"api-guardian/internal/domain/rate_limit/interfaces" // Pastikan interface LimitRepository ada di sini
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// 1. Definisikan Struct (Bukan Function lagi)
type HealthHandler struct {
	Store interfaces.LimitRepository
}

// 2. Constructor untuk Inject Dependency
func NewHealthHandler(store interfaces.LimitRepository) *HealthHandler {
	return &HealthHandler{
		Store: store,
	}
}

// 3. Method Check (Handler-nya)
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	redisStatus := "Connected"

	// Panggil via h.Store (Struct Field), bukan parameter fungsi
	if h.Store != nil {
		if err := h.Store.Ping(context.Background()); err != nil {
			redisStatus = "Disconnected: " + err.Error()
		}
	} else {
		redisStatus = "Not Configured"
	}

	status := map[string]string{
		"system":           "Healthy",
		"redis_connection": redisStatus,
		"circuit_breaker":  "Closed (Normal)",
		"time":             time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
