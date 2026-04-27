package config

import (
	"os"
)

type Config struct {
	ListenAddr string
	DBPath     string
	StaticDir  string // when set, serve frontend SPA from this dir under "/"
}

func Load() (Config, error) {
	return Config{
		ListenAddr: envOr("LISTEN_ADDR", ":8080"),
		DBPath:     envOr("DB_PATH", "data/japanese-site.sqlite"),
		StaticDir:  os.Getenv("STATIC_DIR"),
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
