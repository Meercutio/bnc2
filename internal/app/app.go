package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"example.com/bc-mvp/internal/auth"
	"example.com/bc-mvp/internal/config"
	"example.com/bc-mvp/internal/game"
	"example.com/bc-mvp/internal/httpapi"
	"example.com/bc-mvp/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg config.Config
	log *slog.Logger

	db  *pgxpool.Pool
	rdb *redis.Client

	srv *http.Server
}

type Options struct {
	Static http.Handler // optional; if nil, no frontend is served
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger, opts Options) (*App, error) {
	if log == nil {
		log = slog.Default()
	}

	// --- Postgres ---
	dbpool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool: %w", err)
	}

	// --- Redis ---
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		DB:   cfg.Redis.DB,
	})

	// Quick connectivity checks (fail fast).
	pingCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	if err := dbpool.Ping(pingCtx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	pingErr := rdb.Ping(pingCtx).Err()
	if pingErr != nil {
		dbpool.Close()
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping (%s db=%d): %w", cfg.Redis.Addr, cfg.Redis.DB, pingErr)
	}

	// --- Auth service ---
	authSvc := auth.NewService([]byte(cfg.Auth.Secret))

	// --- Stores ---
	users := store.NewUserStore(dbpool)
	stats := store.NewStatsStore(dbpool)

	authH := &httpapi.AuthHandler{
		Users:    users,
		Stats:    stats,
		Auth:     authSvc,
		TokenTTL: cfg.Auth.TokenTTL,
	}

	// --- Game ---
	persist := game.NewRedisMatchStore(rdb, cfg.Redis.MatchTTL)
	gameCfg := game.Config{RoundDuration: cfg.Game.RoundDuration}
	matchSvc := game.NewMatchService(gameCfg, persist)
	gameSrv := game.NewServer(gameCfg, matchSvc, authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	gameSrv.RegisterRoutes(mux)

	// --- auth routes ---
	mux.HandleFunc("/api/auth/register", authH.Register)
	mux.HandleFunc("/api/auth/login", authH.Login)
	mux.Handle("/api/me", httpapi.AuthMiddleware(authSvc)(http.HandlerFunc(authH.Me)))

	if opts.Static != nil {
		mux.Handle("/", opts.Static)
	}

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}

	return &App{cfg: cfg, log: log, db: dbpool, rdb: rdb, srv: srv}, nil
}

func (a *App) Run(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)

	a.log.Info("http server starting", "addr", a.cfg.HTTP.Addr)

	g.Go(func() error {
		err := a.srv.ListenAndServe()
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	g.Go(func() error {
		<-gctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
		defer cancel()
		a.log.Info("http server shutting down")
		_ = a.srv.Shutdown(shutdownCtx)
		return nil
	})

	err := g.Wait()
	_ = a.Close(context.Background())
	return err
}

func (a *App) Close(ctx context.Context) error {
	// best-effort
	if a.db != nil {
		a.db.Close()
	}
	if a.rdb != nil {
		_ = a.rdb.Close()
	}
	return nil
}
