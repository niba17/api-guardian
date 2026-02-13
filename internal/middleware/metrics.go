package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// --- DEFINISI METRICS ---

// 1. Counter: Menghitung jumlah total request
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_guardian_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)

	// 2. Histogram: Mengukur seberapa cepat (latency) request diproses
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_guardian_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // Bucket standar (0.005s, 0.01s, ..., 10s)
		},
		[]string{"method", "path"},
	)
)

// responseWriterWrapper untuk menangkap Status Code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// PrometheusMiddleware adalah "Petugas Sensus" yang mencatat data
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Bungkus ResponseWriter supaya kita bisa intip status code-nya nanti
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Jalankan request ke middleware berikutnya
		next.ServeHTTP(wrapper, r)

		// --- SETELAH REQUEST SELESAI ---
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapper.statusCode)

		// 1. Catat Durasi
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		// 2. Tambah Counter (sesuai method, path, dan status code)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
	})
}
