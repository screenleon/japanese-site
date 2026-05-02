package config

import (
	"log/slog"
	"os"
)

type Config struct {
	ListenAddr   string
	DBPath       string
	StaticDir    string // when set, serve frontend SPA from this dir under "/"
	ProgressMode string
}

func Load() (Config, error) {
	progressMode := envOr("JS_PROGRESS_STORE", "sqlite")
	switch progressMode {
	case "sqlite", "null":
	default:
		slog.Warn("invalid JS_PROGRESS_STORE, falling back to null", "value", progressMode)
		progressMode = "null"
	}

	return Config{
		ListenAddr:   envOr("LISTEN_ADDR", ":8080"),
		DBPath:       envOr("DB_PATH", "data/japanese-site.sqlite"),
		StaticDir:    os.Getenv("STATIC_DIR"),
		ProgressMode: progressMode,
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
