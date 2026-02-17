package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup mengatur konfigurasi logging global aplikasi
func Setup() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	fileLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     28, // Days
		Compress:   true,
	}

	// MultiWriter: Output ke Konsol (untuk kita pantau) dan ke File (untuk arsip)
	multi := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		fileLogger,
	)

	log.Logger = log.Output(multi)
}
