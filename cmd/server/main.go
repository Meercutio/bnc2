package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/bc-mvp/internal/game"
	"example.com/bc-mvp/internal/httpapi"
	"example.com/bc-mvp/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := env("PORT", "8080")
	addr := ":" + port
	roundDur := envDuration("ROUND_DURATION", 0*time.Second)

	// --- Postgres ---
	dbURL := env("DATABASE_URL", "postgres://bc:bc@localhost:5432/bc?sslmode=disable")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer dbpool.Close()
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	if env("RUN_MIGRATIONS", "false") == "true" {
		runMigrations(dbURL)
	}

	// --- JWT ---
	jwtSecret := []byte(env("JWT_SECRET", "dev-secret-change-me"))
	tokenTTL := envDuration("JWT_TTL", 24*time.Hour)
	game.SetJWTSecret(jwtSecret)

	users := store.NewUserStore(dbpool)
	stats := store.NewStatsStore(dbpool)

	authH := &httpapi.AuthHandler{
		Users:     users,
		Stats:     stats,
		JWTSecret: jwtSecret,
		TokenTTL:  tokenTTL,
	}

	// --- Redis (matches persistence) ---
	redisAddr := env("REDIS_ADDR", "localhost:6379")
	matchTTL := envDuration("MATCH_TTL", 24*time.Hour)

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	persist := game.NewRedisMatchStore(rdb, matchTTL)
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			log.Fatalf("redis close problem: %v", err)
		}
	}(rdb)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	cfg := game.Config{RoundDuration: roundDur}
	matchSvc := game.NewMatchService(cfg, persist)
	srv := game.NewServer(cfg, matchSvc)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// --- game routes ---
	srv.RegisterRoutes(mux)

	// --- auth routes ---
	mux.HandleFunc("/api/auth/register", authH.Register)
	mux.HandleFunc("/api/auth/login", authH.Login)

	// /api/me защищён middleware
	meHandler := http.HandlerFunc(authH.Me)
	mux.Handle("/api/me", httpapi.AuthMiddleware(jwtSecret)(meHandler))

	// --- embedded frontend ---
	h, err := webHandler()
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", h)

	// Запуск сервера в горутине
	go func() {
		log.Printf("listening on %s (round=%s, pg=%s, redis=%s)", addr, roundDur, dbURL, redisAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
	// Обработка SIGTERM/SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
