package main

import (
	"api-guardian/internal/app"
	"api-guardian/internal/config"
	"api-guardian/pkg/logger"
)

func main() {
	logger.Setup()
	cfg := config.Load()

	// Langsung jalankan aplikasi, biarkan package app yang pusing
	if err := app.Run(cfg); err != nil {
		panic(err)
	}
}
