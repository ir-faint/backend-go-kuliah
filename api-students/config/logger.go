package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// parseLevel mengonversi nama tingkat log ke representasi slog.Level.
func parseLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLogger membuat logger terstruktur yang menulis ke dua tujuan sekaligus:
// layar (stdout) agar terlihat saat praktikum, dan file logs/app.log yang
// dirotasi otomatis agar tidak membengkak tanpa batas.
func NewLogger() *slog.Logger {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("gagal membuat folder logs: " + err.Error())
	}

	rotator := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10, // rotasi setiap file mencapai 10 MB
		MaxBackups: 5,  // simpan 5 file lama
		MaxAge:     14, // hapus file yang lebih tua dari 14 hari
		Compress:   true,
	}

	writer := io.MultiWriter(os.Stdout, rotator)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: parseLevel(GetEnv("LOG_LEVEL", "info")),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
