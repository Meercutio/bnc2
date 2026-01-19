package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"example.com/bc-mvp/internal/app"
	"example.com/bc-mvp/internal/config"
	"example.com/bc-mvp/internal/migrate"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	log := newLogger(cfg)

	if cfg.Postgres.RunMigrations {
		if err := migrate.Up(cfg.Postgres.URL, cfg.Postgres.MigrationsDir, log); err != nil {
			log.Error("migrations failed", "err", err)
			os.Exit(1)
		}
	}

	static, err := webHandler()
	if err != nil {
		log.Error("static handler error", "err", err)
		os.Exit(1)
	}

	a, err := app.New(ctx, cfg, log, app.Options{Static: static})
	if err != nil {
		log.Error("app init error", "err", err)
		os.Exit(1)
	}

	log.Info("starting", "env", cfg.Env, "addr", cfg.HTTP.Addr)
	if err := a.Run(ctx); err != nil {
		log.Error("app stopped with error", "err", err)
		os.Exit(1)
	}
}

func newLogger(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if cfg.Log.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
