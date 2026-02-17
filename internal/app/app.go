package app

import (
	"api-guardian/internal/config"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type GuardianApp struct {
	server *http.Server
	cfg    *config.AppConfig
}

func New(cfg *config.AppConfig, handler http.Handler) *GuardianApp {
	return &GuardianApp{
		cfg: cfg,
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: handler,
		},
	}
}

func (a *GuardianApp) Run() {
	go func() {
		log.Info().Str("port", a.cfg.Port).Msg("API Guardian standing guard...")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Stop()
}

func (a *GuardianApp) Stop() {
	log.Info().Msg("Shutting down API Guardian...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server shutdown complete")
}
