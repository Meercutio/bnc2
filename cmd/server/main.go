package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"example.com/bc-mvp/internal/game"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := env("PORT", "8080")                              // Render выставит PORT автоматически
	addr := ":" + port                                       // слушаем 0.0.0.0:<PORT> по умолчанию
	roundDur := envDuration("ROUND_DURATION", 0*time.Second) // 0 => таймер выключен

	cfg := game.Config{
		RoundDuration: roundDur,
	}

	// --- Redis + Persistent match storage ---
	redisAddr := env("REDIS_ADDR", "localhost:6379")
	matchTTL := envDuration("MATCH_TTL", 24*time.Hour)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	persist := game.NewRedisMatchStore(rdb, matchTTL)
	matchSvc := game.NewMatchService(cfg, persist)

	// Server теперь создаётся с match service
	srv := game.NewServer(cfg, matchSvc)

	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	// заменяем статическую раздачу на embedded (перекрываем "/")
	h, err := webHandler()
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", h)

	log.Printf("listening on %s (round duration: %s, redis: %s, match_ttl: %s)", addr, roundDur, redisAddr, matchTTL)
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
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return def
}
