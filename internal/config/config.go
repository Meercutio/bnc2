package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config describes all runtime settings for the server.
//
// Best practice for Go services:
//   - load config once in main
//   - validate
//   - pass further via DI (no global variables)
type Config struct {
	Env string // dev|stage|prod

	Log struct {
		Format string // text|json
	}

	HTTP struct {
		Addr              string
		ReadHeaderTimeout time.Duration
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		IdleTimeout       time.Duration
		ShutdownTimeout   time.Duration
	}

	Postgres struct {
		URL           string
		RunMigrations bool
		MigrationsDir string
	}

	Redis struct {
		Addr     string
		DB       int
		MatchTTL time.Duration
	}

	Auth struct {
		Secret   string
		TokenTTL time.Duration
	}

	Game struct {
		RoundDuration time.Duration
	}
}

func LoadFromEnv() (Config, error) {
	var c Config

	c.Env = envString("APP_ENV", "dev")
	c.Log.Format = envString("LOG_FORMAT", "text")

	port := envString("PORT", "8080")
	c.HTTP.Addr = envString("HTTP_ADDR", ":"+port)
	c.HTTP.ReadHeaderTimeout = envDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	c.HTTP.ReadTimeout = envDuration("HTTP_READ_TIMEOUT", 0)
	c.HTTP.WriteTimeout = envDuration("HTTP_WRITE_TIMEOUT", 0)
	c.HTTP.IdleTimeout = envDuration("HTTP_IDLE_TIMEOUT", 60*time.Second)
	c.HTTP.ShutdownTimeout = envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)

	c.Postgres.URL = envString("DATABASE_URL", "postgres://bc:bc@localhost:5432/bc?sslmode=disable")
	c.Postgres.RunMigrations = envBool("RUN_MIGRATIONS", false)
	c.Postgres.MigrationsDir = envString("MIGRATIONS_DIR", "./db/migrations")

	c.Redis.Addr = envString("REDIS_ADDR", "localhost:6379")
	c.Redis.DB = envInt("REDIS_DB", 0)
	c.Redis.MatchTTL = envDuration("MATCH_TTL", 24*time.Hour)

	c.Auth.Secret = envString("JWT_SECRET", "dev-secret-change-me")
	c.Auth.TokenTTL = envDuration("JWT_TTL", 24*time.Hour)

	c.Game.RoundDuration = envDuration("ROUND_DURATION", 0)

	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if c.HTTP.Addr == "" {
		return errors.New("HTTP addr is empty")
	}
	if c.Postgres.URL == "" {
		return errors.New("DATABASE_URL is empty")
	}
	if c.Redis.Addr == "" {
		return errors.New("REDIS_ADDR is empty")
	}
	if c.Auth.Secret == "" {
		return errors.New("JWT_SECRET is empty")
	}
	if c.Env != "dev" && c.Auth.Secret == "dev-secret-change-me" {
		return fmt.Errorf("refuse to run with default JWT_SECRET in %s", c.Env)
	}
	if c.Log.Format != "text" && c.Log.Format != "json" {
		return fmt.Errorf("unsupported LOG_FORMAT=%q (want text|json)", c.Log.Format)
	}
	return nil
}

func envString(key, def string) string {
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

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}
