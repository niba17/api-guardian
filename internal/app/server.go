package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func startServer(h http.Handler, port string, rdb *redis.Client, geo *geoip2.Reader) error {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msgf("🛡️ API Guardian Standing Guard on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failure")
		}
	}()

	<-stop
	log.Info().Msg("⚠️ Shutting down...")

	// 1. Buat context timeout dulu
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Matikan HTTP Server dulu (tunggu request aktif selesai)
	err := srv.Shutdown(ctx)

	// 3. BARU TUTUP Resource pendukung (Redis & GeoIP)
	// Sekarang aman, karena sudah tidak ada request aktif yang pakai resource ini
	if rdb != nil {
		log.Info().Msg("🔌 Closing Redis connection...")
		rdb.Close()
	}
	if geo != nil {
		log.Info().Msg("🌍 Closing GeoIP database...")
		geo.Close()
	}

	log.Info().Msg("✅ API Guardian has retired for the day. Graceful shutdown complete.")

	return err
}
