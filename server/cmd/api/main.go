package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/screenleon/japanese-site/server/internal/config"
	"github.com/screenleon/japanese-site/server/internal/handlers"
	"github.com/screenleon/japanese-site/server/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		logger.Error("db open failed", "path", cfg.DBPath, "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := store.Migrate(db); err != nil {
		logger.Error("migrate failed", "err", err)
		os.Exit(1)
	}

	var ps store.ProgressStore
	switch cfg.ProgressMode {
	case "sqlite":
		ps = store.NewSQLiteProgressStore(db)
	default:
		ps = store.NullProgressStore{}
	}

	mux := http.NewServeMux()
	handlers.Register(mux, db, ps)
	if err := handlers.RegisterStatic(mux, cfg.StaticDir); err != nil {
		logger.Error("static dir setup failed", "err", err)
		os.Exit(1)
	}
	if cfg.StaticDir != "" {
		logger.Info("serving SPA", "dir", cfg.StaticDir)
	}

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: handlers.SecurityHeaders(mux),
		// Bound every phase of the request lifecycle. Defends against
		// slowloris / slow-body attacks once this is reachable from a
		// network beyond localhost.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server listening", "addr", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "err", err)
		os.Exit(1)
	}
	logger.Info("shutdown complete")
}
