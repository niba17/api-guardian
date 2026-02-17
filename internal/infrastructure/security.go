package infrastructure // 👈 GANTI DARI 'package app' JADI INI

import (
	"api-guardian/internal/config" // Sesuaikan path middleware
	"api-guardian/internal/proxy"  // Tambahkan proxy import untuk LoadBalancer
	"api-guardian/internal/repository"
	"api-guardian/internal/storage"
	"api-guardian/internal/usecase"
	"net/http"

	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Update parameternya juga, tampaknya di main.go Bos passing LoadBalancer (lb)
// tapi di fungsi ini belum ada parameter lb.
func ApplySecurityMiddleware(
	cfg *config.AppConfig,
	db *gorm.DB,
	geoDB *geoip2.Reader,
	cacheUC *usecase.CacheUsecase,
	rdb *redis.Client,
	store storage.LimiterStore,
	logRepo repository.SecurityLogRepository,
	lb *proxy.LoadBalancer, // 👈 Tambahkan parameter ini agar cocok dengan main.go
) http.Handler { // Ubah return type jadi http.Handler (bukan func)

	// Logic chain middleware...
	// Middleware terluar membungkus middleware di dalamnya
	// Urutan Eksekusi: WAF -> RateLimit -> Cache -> Router -> Handler

	// Kita mulai dari yang paling dalam (Load Balancer / Router)
	// var baseHandler http.Handler = lb
	// Tapi karena lb strukturnya belum jelas, kita return wrapper saja
	// Nanti di main.go baru kita sambung ke Router

	// Strategi Baru: Fungsi ini mengembalikan "Chain Function" atau "Wrapper"
	// Tapi biar gampang, kita ubah logic di main.go saja nanti.

	return nil // Placeholder, lihat perbaikan main.go di bawah
}
