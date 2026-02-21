package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// --- DEFINISI METRICS ---

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_guardian_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_guardian_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Regex untuk menyamarkan angka (ID) atau UUID di dalam URL
	// Contoh: /api/menus/123 -> /api/menus/:id
	idRegex = regexp.MustCompile(`/[0-9a-fA-F-]+/?$|/[0-9]+`)
)

// Helper untuk mencegah Prometheus Cardinality Explosion
func normalizePath(path string) string {
	return idRegex.ReplaceAllString(path, "/:id")
}

// 🚀 PERBAIKAN: responseWriterWrapper disamakan standar amannya dengan Logger
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// Tangkap status 200 implisit
func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// PrometheusMiddleware adalah "Petugas Sensus" yang mencatat data
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Bungkus ResponseWriter dengan aman
		wrapper := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default ke 200
		}

		next.ServeHTTP(wrapper, r)

		// --- SETELAH REQUEST SELESAI ---
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapper.statusCode)

		// 🚀 Normalisasi URL sebelum masuk ke Prometheus
		cleanPath := normalizePath(r.URL.Path)

		httpRequestDuration.WithLabelValues(r.Method, cleanPath).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, cleanPath, status).Inc()
	})
}
