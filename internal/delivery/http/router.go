package http

import (
	"api-guardian/internal/delivery/http/handler"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter sekarang menerima Struct HealthHandler, bukan func biasa
func NewRouter(
	authHandler *handler.AuthHandler,
	dashHandler *handler.DashboardHandler,
	healthHandler *handler.HealthHandler,
	authMW func(http.Handler) http.Handler,
	secureChain http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// --- 1. Public Routes ---
	// Langsung ke Handler -> Usecase -> Repo (Sesuai Alur)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("GET /status", healthHandler.Check)
	mux.Handle("/metrics", promhttp.Handler())

	// --- 2. Protected Routes ---
	// (Contoh implementasi dashboard, uncomment jika dashHandler sudah siap)
	if dashHandler != nil {
		mux.Handle("GET /api/dashboard/stats", authMW(http.HandlerFunc(dashHandler.GetDashboardStats)))
		mux.Handle("GET /api/dashboard/logs", authMW(http.HandlerFunc(dashHandler.GetRecentLogs)))
	}

	return mux
}
