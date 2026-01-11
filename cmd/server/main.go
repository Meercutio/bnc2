package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"example.com/bc-mvp/internal/game"
)

func main() {
	addr := env("ADDR", ":8080")
	roundDur := envDuration("ROUND_DURATION", 30*time.Second) // 0 => таймер выключен

	srv := game.NewServer(game.Config{
		RoundDuration: roundDur,
	})

	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	log.Printf("listening on %s (round duration: %s)", addr, roundDur)
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
