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

	// Cleanup
	if rdb != nil {
		rdb.Close()
	}
	if geo != nil {
		geo.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
