package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func startServer(h http.Handler, port string, rdb io.Closer, db io.Closer, geo io.Closer) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Matikan HTTP Server
	err := srv.Shutdown(ctx)

	// 2. Tutup semua Resource (Order: Redis -> DB -> GeoIP)
	if rdb != nil {
		log.Info().Msg("🔌 Closing Redis connection...")
		rdb.Close()
	}

	// 🚀 FIX: Tambahkan logika tutup DB di sini
	if db != nil {
		log.Info().Msg("🐘 Closing PostgreSQL connection...")
		db.Close()
	}

	if geo != nil {
		log.Info().Msg("🌍 Closing GeoIP database...")
		geo.Close()
	}

	log.Info().Msg("✅ API Guardian has retired for the day. Graceful shutdown complete.")

	return err
}
