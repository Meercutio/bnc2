package main

import (
	"context"
	"log"
	"net/http"
	"os"
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

	cfg := game.Config{RoundDuration: roundDur}
	matchSvc := game.NewMatchService(cfg, persist)
	srv := game.NewServer(cfg, matchSvc)

	mux := http.NewServeMux()

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

	log.Printf("listening on %s (round=%s, pg=%s, redis=%s)", addr, roundDur, dbURL, redisAddr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
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
